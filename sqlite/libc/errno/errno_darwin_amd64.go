// Code generated by 'ccgo errno/gen.c -crt-import-path "" -export-defines "" -export-enums "" -export-externs X -export-fields F -export-structs "" -export-typedefs "" -hide _OSSwapInt16,_OSSwapInt32,_OSSwapInt64 -o errno/errno_darwin_amd64.go -pkgname errno', DO NOT EDIT.

package errno

import (
	"math"
	"reflect"
	"sync/atomic"
	"unsafe"
)

var _ = math.Pi
var _ reflect.Kind
var _ atomic.Value
var _ unsafe.Pointer

const (
	E2BIG                                  = 7
	EACCES                                 = 13
	EADDRINUSE                             = 48
	EADDRNOTAVAIL                          = 49
	EAFNOSUPPORT                           = 47
	EAGAIN                                 = 35
	EALREADY                               = 37
	EAUTH                                  = 80
	EBADARCH                               = 86
	EBADEXEC                               = 85
	EBADF                                  = 9
	EBADMACHO                              = 88
	EBADMSG                                = 94
	EBADRPC                                = 72
	EBUSY                                  = 16
	ECANCELED                              = 89
	ECHILD                                 = 10
	ECONNABORTED                           = 53
	ECONNREFUSED                           = 61
	ECONNRESET                             = 54
	EDEADLK                                = 11
	EDESTADDRREQ                           = 39
	EDEVERR                                = 83
	EDOM                                   = 33
	EDQUOT                                 = 69
	EEXIST                                 = 17
	EFAULT                                 = 14
	EFBIG                                  = 27
	EFTYPE                                 = 79
	EHOSTDOWN                              = 64
	EHOSTUNREACH                           = 65
	EIDRM                                  = 90
	EILSEQ                                 = 92
	EINPROGRESS                            = 36
	EINTR                                  = 4
	EINVAL                                 = 22
	EIO                                    = 5
	EISCONN                                = 56
	EISDIR                                 = 21
	ELAST                                  = 106
	ELOOP                                  = 62
	EMFILE                                 = 24
	EMLINK                                 = 31
	EMSGSIZE                               = 40
	EMULTIHOP                              = 95
	ENAMETOOLONG                           = 63
	ENEEDAUTH                              = 81
	ENETDOWN                               = 50
	ENETRESET                              = 52
	ENETUNREACH                            = 51
	ENFILE                                 = 23
	ENOATTR                                = 93
	ENOBUFS                                = 55
	ENODATA                                = 96
	ENODEV                                 = 19
	ENOENT                                 = 2
	ENOEXEC                                = 8
	ENOLCK                                 = 77
	ENOLINK                                = 97
	ENOMEM                                 = 12
	ENOMSG                                 = 91
	ENOPOLICY                              = 103
	ENOPROTOOPT                            = 42
	ENOSPC                                 = 28
	ENOSR                                  = 98
	ENOSTR                                 = 99
	ENOSYS                                 = 78
	ENOTBLK                                = 15
	ENOTCONN                               = 57
	ENOTDIR                                = 20
	ENOTEMPTY                              = 66
	ENOTRECOVERABLE                        = 104
	ENOTSOCK                               = 38
	ENOTSUP                                = 45
	ENOTTY                                 = 25
	ENXIO                                  = 6
	EOPNOTSUPP                             = 102
	EOVERFLOW                              = 84
	EOWNERDEAD                             = 105
	EPERM                                  = 1
	EPFNOSUPPORT                           = 46
	EPIPE                                  = 32
	EPROCLIM                               = 67
	EPROCUNAVAIL                           = 76
	EPROGMISMATCH                          = 75
	EPROGUNAVAIL                           = 74
	EPROTO                                 = 100
	EPROTONOSUPPORT                        = 43
	EPROTOTYPE                             = 41
	EPWROFF                                = 82
	EQFULL                                 = 106
	ERANGE                                 = 34
	EREMOTE                                = 71
	EROFS                                  = 30
	ERPCMISMATCH                           = 73
	ESHLIBVERS                             = 87
	ESHUTDOWN                              = 58
	ESOCKTNOSUPPORT                        = 44
	ESPIPE                                 = 29
	ESRCH                                  = 3
	ESTALE                                 = 70
	ETIME                                  = 101
	ETIMEDOUT                              = 60
	ETOOMANYREFS                           = 59
	ETXTBSY                                = 26
	EUSERS                                 = 68
	EWOULDBLOCK                            = 35
	EXDEV                                  = 18
	X_CDEFS_H_                             = 0
	X_DARWIN_FEATURE_64_BIT_INODE          = 1
	X_DARWIN_FEATURE_ONLY_UNIX_CONFORMANCE = 1
	X_DARWIN_FEATURE_UNIX_CONFORMANCE      = 3
	X_ERRNO_T                              = 0
	X_FILE_OFFSET_BITS                     = 64
	X_LP64                                 = 1
	X_Nonnull                              = 0
	X_Null_unspecified                     = 0
	X_Nullable                             = 0
	X_SYS_ERRNO_H_                         = 0
)

type Ptrdiff_t = int64 /* <builtin>:3:26 */

type Size_t = uint64 /* <builtin>:9:23 */

type Wchar_t = int32 /* <builtin>:15:24 */

type X__int128_t = struct {
	Flo int64
	Fhi int64
} /* <builtin>:21:43 */ // must match utilware/sqlite/libc/mathutil.Int128
type X__uint128_t = struct {
	Flo uint64
	Fhi uint64
} /* <builtin>:22:44 */ // must match utilware/sqlite/libc/mathutil.Int128

type X__builtin_va_list = uintptr /* <builtin>:46:14 */
type X__float128 = float64        /* <builtin>:47:21 */

// Copyright (c) 2000 Apple Computer, Inc. All rights reserved.
//
// @APPLE_LICENSE_HEADER_START@
//
// This file contains Original Code and/or Modifications of Original Code
// as defined in and that are subject to the Apple Public Source License
// Version 2.0 (the 'License'). You may not use this file except in
// compliance with the License. Please obtain a copy of the License at
// http://www.opensource.apple.com/apsl/ and read it before using this
// file.
//
// The Original Code and all software distributed under the License are
// distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
// EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
// INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
// Please see the License for the specific language governing rights and
// limitations under the License.
//
// @APPLE_LICENSE_HEADER_END@
// Copyright (c) 2000-2012 Apple, Inc. All rights reserved.
//
// @APPLE_OSREFERENCE_LICENSE_HEADER_START@
//
// This file contains Original Code and/or Modifications of Original Code
// as defined in and that are subject to the Apple Public Source License
// Version 2.0 (the 'License'). You may not use this file except in
// compliance with the License. The rights granted to you under the License
// may not be used to create, or enable the creation or redistribution of,
// unlawful or unlicensed copies of an Apple operating system, or to
// circumvent, violate, or enable the circumvention or violation of, any
// terms of an Apple operating system software license agreement.
//
// Please obtain a copy of the License at
// http://www.opensource.apple.com/apsl/ and read it before using this file.
//
// The Original Code and all software distributed under the License are
// distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
// EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
// INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
// Please see the License for the specific language governing rights and
// limitations under the License.
//
// @APPLE_OSREFERENCE_LICENSE_HEADER_END@
// Copyright (c) 1995 NeXT Computer, Inc. All Rights Reserved
// Copyright (c) 1982, 1986, 1989, 1993
//	The Regents of the University of California.  All rights reserved.
// (c) UNIX System Laboratories, Inc.
// All or some portions of this file are derived from material licensed
// to the University of California by American Telephone and Telegraph
// Co. or Unix System Laboratories, Inc. and are reproduced herein with
// the permission of UNIX System Laboratories, Inc.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. All advertising materials mentioning features or use of this software
//    must display the following acknowledgement:
//	This product includes software developed by the University of
//	California, Berkeley and its contributors.
// 4. Neither the name of the University nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.
//
//	@(#)errno.h	8.5 (Berkeley) 1/21/94

// Copyright (c) 2000-2018 Apple Inc. All rights reserved.
//
// @APPLE_OSREFERENCE_LICENSE_HEADER_START@
//
// This file contains Original Code and/or Modifications of Original Code
// as defined in and that are subject to the Apple Public Source License
// Version 2.0 (the 'License'). You may not use this file except in
// compliance with the License. The rights granted to you under the License
// may not be used to create, or enable the creation or redistribution of,
// unlawful or unlicensed copies of an Apple operating system, or to
// circumvent, violate, or enable the circumvention or violation of, any
// terms of an Apple operating system software license agreement.
//
// Please obtain a copy of the License at
// http://www.opensource.apple.com/apsl/ and read it before using this file.
//
// The Original Code and all software distributed under the License are
// distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
// EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
// INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
// Please see the License for the specific language governing rights and
// limitations under the License.
//
// @APPLE_OSREFERENCE_LICENSE_HEADER_END@
// Copyright 1995 NeXT Computer, Inc. All rights reserved.
// Copyright (c) 1991, 1993
//	The Regents of the University of California.  All rights reserved.
//
// This code is derived from software contributed to Berkeley by
// Berkeley Software Design, Inc.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. All advertising materials mentioning features or use of this software
//    must display the following acknowledgement:
//	This product includes software developed by the University of
//	California, Berkeley and its contributors.
// 4. Neither the name of the University nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.
//
//	@(#)cdefs.h	8.8 (Berkeley) 1/9/95

// This SDK is designed to work with clang and specific versions of
// gcc >= 4.0 with Apple's patch sets

// Compatibility with compilers and environments that don't support compiler
// feature checking function-like macros.

// The __CONCAT macro is used to concatenate parts of symbol names, e.g.
// with "#define OLD(foo) __CONCAT(old,foo)", OLD(foo) produces oldfoo.
// The __CONCAT macro is a bit tricky -- make sure you don't put spaces
// in between its arguments.  __CONCAT can also concatenate double-quoted
// strings produced by the __STRING macro, but this only works with ANSI C.

// __unused denotes variables and functions that may not be used, preventing
// the compiler from warning about it if not used.

// __used forces variables and functions to be included even if it appears
// to the compiler that they are not used (and would thust be discarded).

// __cold marks code used for debugging or that is rarely taken
// and tells the compiler to optimize for size and outline code.

// __deprecated causes the compiler to produce a warning when encountering
// code using the deprecated functionality.
// __deprecated_msg() does the same, and compilers that support it will print
// a message along with the deprecation warning.
// This may require turning on such warning with the -Wdeprecated flag.
// __deprecated_enum_msg() should be used on enums, and compilers that support
// it will print the deprecation warning.
// __kpi_deprecated() specifically indicates deprecation of kernel programming
// interfaces in Kernel.framework used by KEXTs.

// __unavailable causes the compiler to error out when encountering
// code using the tagged function of variable.

// Delete pseudo-keywords wherever they are not available or needed.

// We use `__restrict' as a way to define the `restrict' type qualifier
// without disturbing older software that is unaware of C99 keywords.

// Compatibility with compilers and environments that don't support the
// nullability feature.

// __disable_tail_calls causes the compiler to not perform tail call
// optimization inside the marked function.

// __not_tail_called causes the compiler to prevent tail call optimization
// on statically bound calls to the function.  It has no effect on indirect
// calls.  Virtual functions, objective-c methods, and functions marked as
// "always_inline" cannot be marked as __not_tail_called.

// __result_use_check warns callers of a function that not using the function
// return value is a bug, i.e. dismissing malloc() return value results in a
// memory leak.

// __swift_unavailable causes the compiler to mark a symbol as specifically
// unavailable in Swift, regardless of any other availability in C.

// __abortlike is the attribute to put on functions like abort() that are
// typically used to mark assertions. These optimize the codegen
// for outlining while still maintaining debugability.

// Declaring inline functions within headers is error-prone due to differences
// across various versions of the C language and extensions.  __header_inline
// can be used to declare inline functions within system headers.  In cases
// where you want to force inlining instead of letting the compiler make
// the decision, you can use __header_always_inline.
//
// Be aware that using inline for functions which compilers may also provide
// builtins can behave differently under various compilers.  If you intend to
// provide an inline version of such a function, you may want to use a macro
// instead.
//
// The check for !__GNUC__ || __clang__ is because gcc doesn't correctly
// support c99 inline in some cases:
// http://gcc.gnu.org/bugzilla/show_bug.cgi?id=55965

// Compiler-dependent macros that bracket portions of code where the
// "-Wunreachable-code" warning should be ignored. Please use sparingly.

// Compiler-dependent macros to declare that functions take printf-like
// or scanf-like arguments.  They are null except for versions of gcc
// that are known to support the features properly.  Functions declared
// with these attributes will cause compilation warnings if there is a
// mismatch between the format string and subsequent function parameter
// types.

// Source compatibility only, ID string not emitted in object file

// __alloc_size can be used to label function arguments that represent the
// size of memory that the function allocates and returns. The one-argument
// form labels a single argument that gives the allocation size (where the
// arguments are numbered from 1):
//
// void	*malloc(size_t __size) __alloc_size(1);
//
// The two-argument form handles the case where the size is calculated as the
// product of two arguments:
//
// void	*calloc(size_t __count, size_t __size) __alloc_size(1,2);

// COMPILATION ENVIRONMENTS -- see compat(5) for additional detail
//
// DEFAULT	By default newly complied code will get POSIX APIs plus
//		Apple API extensions in scope.
//
//		Most users will use this compilation environment to avoid
//		behavioral differences between 32 and 64 bit code.
//
// LEGACY	Defining _NONSTD_SOURCE will get pre-POSIX APIs plus Apple
//		API extensions in scope.
//
//		This is generally equivalent to the Tiger release compilation
//		environment, except that it cannot be applied to 64 bit code;
//		its use is discouraged.
//
//		We expect this environment to be deprecated in the future.
//
// STRICT	Defining _POSIX_C_SOURCE or _XOPEN_SOURCE restricts the
//		available APIs to exactly the set of APIs defined by the
//		corresponding standard, based on the value defined.
//
//		A correct, portable definition for _POSIX_C_SOURCE is 200112L.
//		A correct, portable definition for _XOPEN_SOURCE is 600L.
//
//		Apple API extensions are not visible in this environment,
//		which can cause Apple specific code to fail to compile,
//		or behave incorrectly if prototypes are not in scope or
//		warnings about missing prototypes are not enabled or ignored.
//
// In any compilation environment, for correct symbol resolution to occur,
// function prototypes must be in scope.  It is recommended that all Apple
// tools users add either the "-Wall" or "-Wimplicit-function-declaration"
// compiler flags to their projects to be warned when a function is being
// used without a prototype in scope.

// These settings are particular to each product.
// Platform: MacOSX
// #undef __DARWIN_ONLY_UNIX_CONFORMANCE (automatically set for 64-bit)

// The __DARWIN_ALIAS macros are used to do symbol renaming; they allow
// legacy code to use the old symbol, thus maintaining binary compatibility
// while new code can use a standards compliant version of the same function.
//
// __DARWIN_ALIAS is used by itself if the function signature has not
// changed, it is used along with a #ifdef check for __DARWIN_UNIX03
// if the signature has changed.  Because the __LP64__ environment
// only supports UNIX03 semantics it causes __DARWIN_UNIX03 to be
// defined, but causes __DARWIN_ALIAS to do no symbol mangling.
//
// As a special case, when XCode is used to target a specific version of the
// OS, the manifest constant __ENVIRONMENT_MAC_OS_X_VERSION_MIN_REQUIRED__
// will be defined by the compiler, with the digits representing major version
// time 100 + minor version times 10 (e.g. 10.5 := 1050).  If we are targeting
// pre-10.5, and it is the default compilation environment, revert the
// compilation environment to pre-__DARWIN_UNIX03.

// symbol suffixes used for symbol versioning

// symbol versioning macros

// symbol release macros
// Copyright (c) 2010 Apple Inc. All rights reserved.
//
// @APPLE_OSREFERENCE_LICENSE_HEADER_START@
//
// This file contains Original Code and/or Modifications of Original Code
// as defined in and that are subject to the Apple Public Source License
// Version 2.0 (the 'License'). You may not use this file except in
// compliance with the License. The rights granted to you under the License
// may not be used to create, or enable the creation or redistribution of,
// unlawful or unlicensed copies of an Apple operating system, or to
// circumvent, violate, or enable the circumvention or violation of, any
// terms of an Apple operating system software license agreement.
//
// Please obtain a copy of the License at
// http://www.opensource.apple.com/apsl/ and read it before using this file.
//
// The Original Code and all software distributed under the License are
// distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
// EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
// INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
// Please see the License for the specific language governing rights and
// limitations under the License.
//
// @APPLE_OSREFERENCE_LICENSE_HEADER_END@

// POSIX.1 requires that the macros we test be defined before any standard
// header file is included.  This permits us to convert values for feature
// testing, as necessary, using only _POSIX_C_SOURCE.
//
// Here's a quick run-down of the versions:
//  defined(_POSIX_SOURCE)		1003.1-1988
//  _POSIX_C_SOURCE == 1L		1003.1-1990
//  _POSIX_C_SOURCE == 2L		1003.2-1992 C Language Binding Option
//  _POSIX_C_SOURCE == 199309L		1003.1b-1993
//  _POSIX_C_SOURCE == 199506L		1003.1c-1995, 1003.1i-1995,
//					and the omnibus ISO/IEC 9945-1: 1996
//  _POSIX_C_SOURCE == 200112L		1003.1-2001
//  _POSIX_C_SOURCE == 200809L		1003.1-2008
//
// In addition, the X/Open Portability Guide, which is now the Single UNIX
// Specification, defines a feature-test macro which indicates the version of
// that specification, and which subsumes _POSIX_C_SOURCE.

// Deal with IEEE Std. 1003.1-1990, in which _POSIX_C_SOURCE == 1L.

// Deal with IEEE Std. 1003.2-1992, in which _POSIX_C_SOURCE == 2L.

// Deal with various X/Open Portability Guides and Single UNIX Spec.

// Deal with all versions of POSIX.  The ordering relative to the tests above is
// important.

// POSIX C deprecation macros
// Copyright (c) 2010 Apple Inc. All rights reserved.
//
// @APPLE_OSREFERENCE_LICENSE_HEADER_START@
//
// This file contains Original Code and/or Modifications of Original Code
// as defined in and that are subject to the Apple Public Source License
// Version 2.0 (the 'License'). You may not use this file except in
// compliance with the License. The rights granted to you under the License
// may not be used to create, or enable the creation or redistribution of,
// unlawful or unlicensed copies of an Apple operating system, or to
// circumvent, violate, or enable the circumvention or violation of, any
// terms of an Apple operating system software license agreement.
//
// Please obtain a copy of the License at
// http://www.opensource.apple.com/apsl/ and read it before using this file.
//
// The Original Code and all software distributed under the License are
// distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
// EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
// INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
// Please see the License for the specific language governing rights and
// limitations under the License.
//
// @APPLE_OSREFERENCE_LICENSE_HEADER_END@

// Set a single macro which will always be defined and can be used to determine
// the appropriate namespace.  For POSIX, these values will correspond to
// _POSIX_C_SOURCE value.  Currently there are two additional levels corresponding
// to ANSI (_ANSI_SOURCE) and Darwin extensions (_DARWIN_C_SOURCE)

// If the developer has neither requested a strict language mode nor a version
// of POSIX, turn on functionality provided by __STDC_WANT_LIB_EXT1__ as part
// of __DARWIN_C_FULL.

// long long is not supported in c89 (__STRICT_ANSI__), but g++ -ansi and
// c99 still want long longs.  While not perfect, we allow long longs for
// g++.

// ****************************************
//
//  Public darwin-specific feature macros
//

// _DARWIN_FEATURE_64_BIT_INODE indicates that the ino_t type is 64-bit, and
// structures modified for 64-bit inodes (like struct stat) will be used.

// _DARWIN_FEATURE_64_ONLY_BIT_INODE indicates that the ino_t type may only
// be 64-bit; there is no support for 32-bit ino_t when this macro is defined
// (and non-zero).  There is no struct stat64 either, as the regular
// struct stat will already be the 64-bit version.

// _DARWIN_FEATURE_ONLY_VERS_1050 indicates that only those APIs updated
// in 10.5 exists; no pre-10.5 variants are available.

// _DARWIN_FEATURE_ONLY_UNIX_CONFORMANCE indicates only UNIX conforming API
// are available (the legacy BSD APIs are not available)

// _DARWIN_FEATURE_UNIX_CONFORMANCE indicates whether UNIX conformance is on,
// and specifies the conformance level (3 is SUSv3)

// This macro casts away the qualifier from the variable
//
// Note: use at your own risk, removing qualifiers can result in
// catastrophic run-time failures.

// __XNU_PRIVATE_EXTERN is a linkage decoration indicating that a symbol can be
// used from other compilation units, but not other libraries or executables.

// Architecture validation for current SDK

// Similar to OS_ENUM/OS_CLOSED_ENUM/OS_OPTIONS/OS_CLOSED_OPTIONS
//
// This provides more advanced type checking on compilers supporting
// the proper extensions, even in C.

// Copyright (c) 2003-2012 Apple Inc. All rights reserved.
//
// @APPLE_OSREFERENCE_LICENSE_HEADER_START@
//
// This file contains Original Code and/or Modifications of Original Code
// as defined in and that are subject to the Apple Public Source License
// Version 2.0 (the 'License'). You may not use this file except in
// compliance with the License. The rights granted to you under the License
// may not be used to create, or enable the creation or redistribution of,
// unlawful or unlicensed copies of an Apple operating system, or to
// circumvent, violate, or enable the circumvention or violation of, any
// terms of an Apple operating system software license agreement.
//
// Please obtain a copy of the License at
// http://www.opensource.apple.com/apsl/ and read it before using this file.
//
// The Original Code and all software distributed under the License are
// distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
// EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
// INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
// Please see the License for the specific language governing rights and
// limitations under the License.
//
// @APPLE_OSREFERENCE_LICENSE_HEADER_END@
type Errno_t = int32 /* _errno_t.h:30:32 */

// Error codes

// 11 was EAGAIN

// math software

// non-blocking and interrupt i/o

// ipc/network software -- argument errors

// ipc/network software -- operational errors

// should be rearranged

// quotas & mush

// Network File System

// Intelligent device errors

// Program loading errors

// This value is only discrete when compiling __DARWIN_UNIX03, or KERNEL

var _ int8 /* gen.c:2:13: */
