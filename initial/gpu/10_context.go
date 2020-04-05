package gpu

//	#cgo pkg-config: OpenCL
//	#cgo CFLAGS: -Wall
//	#include <errno.h>
//	#include <stdlib.h>
//	#include <CL/cl.h>
//
//	static void create_context(void **_ctx, void **_device_id) {
//		cl_context *ctx = (cl_context *) malloc(sizeof(cl_context));
//		if (ctx == NULL) {
//			errno = ENOBUFS;
//			return;
//		}
//		cl_device_id *device_id = (cl_device_id *) malloc(sizeof(cl_device_id));
//		if (device_id == NULL) {
//			free(ctx);
//			errno = ENOBUFS;
//			return;
//		}
//		*_ctx = ctx;
//		*_device_id = device_id;
//		cl_platform_id platform_id[1];
//		cl_uint num_platforms;
//		cl_int ret;
//		if ((ret = clGetPlatformIDs(1, platform_id, &num_platforms)) != CL_SUCCESS) {
//			free(ctx);
//			free(device_id);
//			errno = ret;
//			return;
//		}
//		cl_uint num_devices;
//		if ((ret = clGetDeviceIDs(platform_id[0], CL_DEVICE_TYPE_GPU, 1, device_id, &num_devices)) != CL_SUCCESS) {
//			free(ctx);
//			free(device_id);
//			errno = ret;
//			return;
//		}
//		*ctx = clCreateContext(NULL, 1, device_id, NULL, NULL, &ret);
//		if (ret != CL_SUCCESS) {
//			free(ctx);
//			free(device_id);
//			errno = ret;
//			return;
//		}
//		errno = 0;
//	}
//
//	static void cleanup_context(void *_ctx, void *_device_id) {
//		cl_context *ctx = (cl_context *) _ctx;
//		cl_device_id *device_id = (cl_device_id *) _device_id;
//		clReleaseContext(*ctx);
//		free(ctx);
//		free(device_id);
//		errno = 0;
//	}
import "C"
import (
	"unsafe"
)

var (
	context  unsafe.Pointer
	deviceID unsafe.Pointer
)

func init() {
	if _, e := C.create_context(&context, &deviceID); e != nil {
		panic(convertError(e))
	}
	cleanupFuncs = append(cleanupFuncs, func() {
		C.cleanup_context(context, deviceID)
	})
}
