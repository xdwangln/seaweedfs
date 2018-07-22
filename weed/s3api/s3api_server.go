package s3api

import (
	_ "github.com/chrislusf/seaweedfs/weed/filer2/cassandra"
	_ "github.com/chrislusf/seaweedfs/weed/filer2/leveldb"
	_ "github.com/chrislusf/seaweedfs/weed/filer2/memdb"
	_ "github.com/chrislusf/seaweedfs/weed/filer2/mysql"
	_ "github.com/chrislusf/seaweedfs/weed/filer2/postgres"
	_ "github.com/chrislusf/seaweedfs/weed/filer2/redis"
	"github.com/gorilla/mux"
	"net/http"
)

type S3ApiServerOption struct {
	Filer            string
	FilerGrpcAddress string
	DomainName       string
	BucketsPath      string
}

type S3ApiServer struct {
	option *S3ApiServerOption
}

func NewS3ApiServer(router *mux.Router, option *S3ApiServerOption) (s3ApiServer *S3ApiServer, err error) {
	s3ApiServer = &S3ApiServer{
		option: option,
	}

	s3ApiServer.registerRouter(router)

	return s3ApiServer, nil
}

func (s3a *S3ApiServer) registerRouter(router *mux.Router) {
	// API Router
	apiRouter := router.PathPrefix("/").Subrouter()
	var routers []*mux.Router
	if s3a.option.DomainName != "" {
		routers = append(routers, apiRouter.Host("{bucket:.+}."+s3a.option.DomainName).Subrouter())
	}
	routers = append(routers, apiRouter.PathPrefix("/{bucket}").Subrouter())

	for _, bucket := range routers {

		// PutObject
		bucket.Methods("PUT").Path("/{object:.+}").HandlerFunc(s3a.PutObjectHandler)
		// GetObject
		bucket.Methods("GET").Path("/{object:.+}").HandlerFunc(s3a.GetObjectHandler)
		// HeadObject
		bucket.Methods("HEAD").Path("/{object:.+}").HandlerFunc(s3a.HeadObjectHandler)
		// DeleteObject
		bucket.Methods("DELETE").Path("/{object:.+}").HandlerFunc(s3a.DeleteObjectHandler)

		// PutBucket
		bucket.Methods("PUT").HandlerFunc(s3a.PutBucketHandler)
		// DeleteBucket
		bucket.Methods("DELETE").HandlerFunc(s3a.DeleteBucketHandler)
		// HeadBucket
		bucket.Methods("HEAD").HandlerFunc(s3a.HeadBucketHandler)

		// ListObjectsV1 (Legacy)
		bucket.Methods("GET").HandlerFunc(s3a.ListObjectsV1Handler)

		/*
			// CopyObject
			bucket.Methods("PUT").Path("/{object:.+}").HeadersRegexp("X-Amz-Copy-Source", ".*?(\\/|%2F).*?").HandlerFunc(s3a.CopyObjectHandler)

			// CopyObjectPart
			bucket.Methods("PUT").Path("/{object:.+}").HeadersRegexp("X-Amz-Copy-Source", ".*?(\\/|%2F).*?").HandlerFunc(s3a.CopyObjectPartHandler).Queries("partNumber", "{partNumber:[0-9]+}", "uploadId", "{uploadId:.*}")
			// PutObjectPart
			bucket.Methods("PUT").Path("/{object:.+}").HandlerFunc(s3a.PutObjectPartHandler).Queries("partNumber", "{partNumber:[0-9]+}", "uploadId", "{uploadId:.*}")
			// ListObjectPxarts
			bucket.Methods("GET").Path("/{object:.+}").HandlerFunc(s3a.ListObjectPartsHandler).Queries("uploadId", "{uploadId:.*}")
			// CompleteMultipartUpload
			bucket.Methods("POST").Path("/{object:.+}").HandlerFunc(s3a.CompleteMultipartUploadHandler).Queries("uploadId", "{uploadId:.*}")
			// NewMultipartUpload
			bucket.Methods("POST").Path("/{object:.+}").HandlerFunc(s3a.NewMultipartUploadHandler).Queries("uploads", "")
			// AbortMultipartUpload
			bucket.Methods("DELETE").Path("/{object:.+}").HandlerFunc(s3a.AbortMultipartUploadHandler).Queries("uploadId", "{uploadId:.*}")

			// ListMultipartUploads
			bucket.Methods("GET").HandlerFunc(s3a.ListMultipartUploadsHandler).Queries("uploads", "")
			// ListObjectsV2
			bucket.Methods("GET").HandlerFunc(s3a.ListObjectsV2Handler).Queries("list-type", "2")
			// ListObjectsV1 (Legacy)
			bucket.Methods("GET").HandlerFunc(s3a.ListObjectsV1Handler)
			// DeleteMultipleObjects
			bucket.Methods("POST").HandlerFunc(s3a.DeleteMultipleObjectsHandler).Queries("delete", "")

			// not implemented
			// GetBucketLocation
			bucket.Methods("GET").HandlerFunc(s3a.GetBucketLocationHandler).Queries("location", "")
			// GetBucketPolicy
			bucket.Methods("GET").HandlerFunc(s3a.GetBucketPolicyHandler).Queries("policy", "")
			// GetObjectACL
			bucket.Methods("GET").Path("/{object:.+}").HandlerFunc(s3a.GetObjectACLHandler).Queries("acl", "")
			// GetBucketACL
			bucket.Methods("GET").HandlerFunc(s3a.GetBucketACLHandler).Queries("acl", "")
			// PutBucketPolicy
			bucket.Methods("PUT").HandlerFunc(s3a.PutBucketPolicyHandler).Queries("policy", "")
			// DeleteBucketPolicy
			bucket.Methods("DELETE").HandlerFunc(s3a.DeleteBucketPolicyHandler).Queries("policy", "")
			// PostPolicy
			bucket.Methods("POST").HeadersRegexp("Content-Type", "multipart/form-data*").HandlerFunc(s3a.PostPolicyBucketHandler)
		*/

	}

	// ListBuckets
	apiRouter.Methods("GET").Path("/").HandlerFunc(s3a.ListBucketsHandler)

	// NotFound
	apiRouter.NotFoundHandler = http.HandlerFunc(notFoundHandler)

}