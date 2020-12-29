package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"

	"net/url"

	"github.com/brxie/ebazarek-backend/controller/image"
	"github.com/brxie/ebazarek-backend/utils"
	"github.com/brxie/ebazarek-backend/utils/ilog"
)

type UploadResponse struct {
	Path            string `json:"path"`
	ThumbnailPath   string `json:"thumbnailPath"`
	ThumbnailWidth  int    `json:"thumbnailWidth"`
	ThumbnailHeight int    `json:"thumbnailHeight"`
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	imageID, err := GetUrlParamValue(r, "imageID")
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		return
	}

	imageID, err = url.QueryUnescape(imageID)
	if err != nil {
		utils.WriteMessageResponse(&w, http.StatusBadRequest,
			"Can't unescape image name parameter")
		return
	}

	img, err := image.GetImage(imageID)
	defer img.Close()
	if err != nil {
		if os.IsNotExist(err) {
			utils.WriteMessageResponse(&w, http.StatusNotFound,
				http.StatusText(http.StatusNotFound))
			return
		}
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		return
	}

	w.WriteHeader(http.StatusOK)
	io.Copy(w, img)
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		utils.WriteMessageResponse(&w, http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest)+err.Error())
		return
	}

	multiReader := multipart.NewReader(r.Body, params["boundary"])

	imageBytes := make(map[string][]byte)
	for {
		part, err := multiReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			ilog.Error(err)
			utils.WriteMessageResponse(&w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError))
			return
		}

		bytes, err := ioutil.ReadAll(part)
		if err != nil {
			log.Fatal(err)
		}
		imageBytes[part.FileName()] = append(imageBytes[part.FileName()], bytes...)
	}

	var response UploadResponse
	for fileName, bytes := range imageBytes {
		result, err := image.UploadImage(fileName, bytes)

		response = UploadResponse{
			Path:            result.Path,
			ThumbnailPath:   result.ThumbnailPath,
			ThumbnailWidth:  result.ThumbnailWidth,
			ThumbnailHeight: result.ThumbnailHeight,
		}

		if err != nil {
			ilog.Error(err)
			utils.WriteMessageResponse(&w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError))
			return
		}
		// only single image in one request is supported. Other images
		// will be ignored.
		break
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
