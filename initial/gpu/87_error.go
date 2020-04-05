package gpu

//	#cgo pkg-config: OpenCL
//	#cgo CFLAGS: -Wall
//	#include <errno.h>
//	#include <string.h>
//	#include <CL/cl.h>
//
//	#define ERR(x) case x: return #x;
//
//	static const char *opencl_strerror(cl_int error) {
//		switch (error) {
//			ERR(CL_DEVICE_NOT_FOUND)
//			ERR(CL_DEVICE_NOT_AVAILABLE)
//			ERR(CL_COMPILER_NOT_AVAILABLE)
//			ERR(CL_MEM_OBJECT_ALLOCATION_FAILURE)
//			ERR(CL_OUT_OF_RESOURCES)
//			ERR(CL_OUT_OF_HOST_MEMORY)
//			ERR(CL_PROFILING_INFO_NOT_AVAILABLE)
//			ERR(CL_MEM_COPY_OVERLAP)
//			ERR(CL_IMAGE_FORMAT_MISMATCH)
//			ERR(CL_IMAGE_FORMAT_NOT_SUPPORTED)
//			ERR(CL_BUILD_PROGRAM_FAILURE)
//			ERR(CL_MAP_FAILURE)
//			ERR(CL_MISALIGNED_SUB_BUFFER_OFFSET)
//			ERR(CL_COMPILE_PROGRAM_FAILURE)
//			ERR(CL_LINKER_NOT_AVAILABLE)
//			ERR(CL_LINK_PROGRAM_FAILURE)
//			ERR(CL_DEVICE_PARTITION_FAILED)
//			ERR(CL_KERNEL_ARG_INFO_NOT_AVAILABLE)
//			ERR(CL_INVALID_VALUE)
//			ERR(CL_INVALID_DEVICE_TYPE)
//			ERR(CL_INVALID_PLATFORM)
//			ERR(CL_INVALID_DEVICE)
//			ERR(CL_INVALID_CONTEXT)
//			ERR(CL_INVALID_QUEUE_PROPERTIES)
//			ERR(CL_INVALID_COMMAND_QUEUE)
//			ERR(CL_INVALID_HOST_PTR)
//			ERR(CL_INVALID_MEM_OBJECT)
//			ERR(CL_INVALID_IMAGE_FORMAT_DESCRIPTOR)
//			ERR(CL_INVALID_IMAGE_SIZE)
//			ERR(CL_INVALID_SAMPLER)
//			ERR(CL_INVALID_BINARY)
//			ERR(CL_INVALID_BUILD_OPTIONS)
//			ERR(CL_INVALID_PROGRAM)
//			ERR(CL_INVALID_PROGRAM_EXECUTABLE)
//			ERR(CL_INVALID_KERNEL_NAME)
//			ERR(CL_INVALID_KERNEL_DEFINITION)
//			ERR(CL_INVALID_KERNEL)
//			ERR(CL_INVALID_ARG_INDEX)
//			ERR(CL_INVALID_ARG_VALUE)
//			ERR(CL_INVALID_ARG_SIZE)
//			ERR(CL_INVALID_KERNEL_ARGS)
//			ERR(CL_INVALID_WORK_DIMENSION)
//			ERR(CL_INVALID_WORK_GROUP_SIZE)
//			ERR(CL_INVALID_WORK_ITEM_SIZE)
//			ERR(CL_INVALID_GLOBAL_OFFSET)
//			ERR(CL_INVALID_EVENT_WAIT_LIST)
//			ERR(CL_INVALID_EVENT)
//			ERR(CL_INVALID_OPERATION)
//			ERR(CL_INVALID_GL_OBJECT)
//			ERR(CL_INVALID_BUFFER_SIZE)
//			ERR(CL_INVALID_MIP_LEVEL)
//			ERR(CL_INVALID_GLOBAL_WORK_SIZE)
//			ERR(CL_INVALID_PROPERTY)
//			ERR(CL_INVALID_IMAGE_DESCRIPTOR)
//			ERR(CL_INVALID_COMPILER_OPTIONS)
//			ERR(CL_INVALID_LINKER_OPTIONS)
//			ERR(CL_INVALID_DEVICE_PARTITION_COUNT)
//			default:
//				return "Unknown error";
//		}
//	}
import "C"
import (
	"fmt"
	"syscall"
)

func convertError(e error) error {
	if e == nil {
		return nil
	}
	errno := int(e.(syscall.Errno))
	if errno < 0 {
		return fmt.Errorf("OpenCL Error #%d: %s", errno, C.GoString(C.opencl_strerror(C.int(errno))))
	} else if errno > 0 {
		return fmt.Errorf("System Error #%d: %s", errno, C.GoString(C.strerror(C.int(errno))))
	} else {
		return nil
	}
}
