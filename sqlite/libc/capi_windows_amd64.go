// Code generated by 'go generate' - DO NOT EDIT.

package libc // import "utilware/sqlite/libc"

var CAPI = map[string]struct{}{
	"AccessCheck":                       {},
	"AddAccessDeniedAce":                {},
	"AddAce":                            {},
	"AreFileApisANSI":                   {},
	"BuildCommDCBW":                     {},
	"CancelSynchronousIo":               {},
	"CharLowerW":                        {},
	"ClearCommError":                    {},
	"CloseHandle":                       {},
	"CopyFileW":                         {},
	"CreateDirectoryW":                  {},
	"CreateEventA":                      {},
	"CreateEventW":                      {},
	"CreateFileA":                       {},
	"CreateFileMappingA":                {},
	"CreateFileMappingW":                {},
	"CreateFileW":                       {},
	"CreateHardLinkW":                   {},
	"CreateMutexW":                      {},
	"CreatePipe":                        {},
	"CreateProcessA":                    {},
	"CreateProcessW":                    {},
	"CreateThread":                      {},
	"CreateWindowExW":                   {},
	"DdeAbandonTransaction":             {},
	"DdeAccessData":                     {},
	"DdeClientTransaction":              {},
	"DdeConnect":                        {},
	"DdeCreateDataHandle":               {},
	"DdeCreateStringHandleW":            {},
	"DdeDisconnect":                     {},
	"DdeFreeDataHandle":                 {},
	"DdeFreeStringHandle":               {},
	"DdeGetData":                        {},
	"DdeGetLastError":                   {},
	"DdeInitializeW":                    {},
	"DdeNameService":                    {},
	"DdeQueryStringW":                   {},
	"DdeUnaccessData":                   {},
	"DdeUninitialize":                   {},
	"DebugBreak":                        {},
	"DefWindowProcW":                    {},
	"DeleteCriticalSection":             {},
	"DeleteFileA":                       {},
	"DeleteFileW":                       {},
	"DestroyWindow":                     {},
	"DeviceIoControl":                   {},
	"DispatchMessageW":                  {},
	"DuplicateHandle":                   {},
	"EnterCriticalSection":              {},
	"EnumWindows":                       {},
	"EqualSid":                          {},
	"EscapeCommFunction":                {},
	"ExitProcess":                       {},
	"FindClose":                         {},
	"FindFirstFileExW":                  {},
	"FindFirstFileW":                    {},
	"FindNextFileW":                     {},
	"FlushFileBuffers":                  {},
	"FlushViewOfFile":                   {},
	"FormatMessageA":                    {},
	"FormatMessageW":                    {},
	"FreeLibrary":                       {},
	"GetACP":                            {},
	"GetAce":                            {},
	"GetAclInformation":                 {},
	"GetCommModemStatus":                {},
	"GetCommState":                      {},
	"GetCommandLineW":                   {},
	"GetComputerNameW":                  {},
	"GetConsoleCP":                      {},
	"GetConsoleMode":                    {},
	"GetConsoleScreenBufferInfo":        {},
	"GetCurrentDirectoryW":              {},
	"GetCurrentProcess":                 {},
	"GetCurrentProcessId":               {},
	"GetCurrentThread":                  {},
	"GetCurrentThreadId":                {},
	"GetDiskFreeSpaceA":                 {},
	"GetDiskFreeSpaceW":                 {},
	"GetEnvironmentVariableA":           {},
	"GetEnvironmentVariableW":           {},
	"GetExitCodeProcess":                {},
	"GetExitCodeThread":                 {},
	"GetFileAttributesA":                {},
	"GetFileAttributesExW":              {},
	"GetFileAttributesW":                {},
	"GetFileInformationByHandle":        {},
	"GetFileSecurityA":                  {},
	"GetFileSecurityW":                  {},
	"GetFileSize":                       {},
	"GetFileType":                       {},
	"GetFullPathNameA":                  {},
	"GetFullPathNameW":                  {},
	"GetLastError":                      {},
	"GetLengthSid":                      {},
	"GetLogicalDriveStringsA":           {},
	"GetMessageW":                       {},
	"GetModuleFileNameA":                {},
	"GetModuleFileNameW":                {},
	"GetModuleHandleW":                  {},
	"GetNamedSecurityInfoW":             {},
	"GetOverlappedResult":               {},
	"GetPrivateProfileStringA":          {},
	"GetProcAddress":                    {},
	"GetProcessHeap":                    {},
	"GetProfilesDirectoryW":             {},
	"GetSecurityDescriptorDacl":         {},
	"GetSecurityDescriptorOwner":        {},
	"GetShortPathNameW":                 {},
	"GetSidIdentifierAuthority":         {},
	"GetSidLengthRequired":              {},
	"GetSidSubAuthority":                {},
	"GetStdHandle":                      {},
	"GetSystemInfo":                     {},
	"GetSystemTime":                     {},
	"GetSystemTimeAsFileTime":           {},
	"GetTempFileNameW":                  {},
	"GetTempPathA":                      {},
	"GetTempPathW":                      {},
	"GetTickCount":                      {},
	"GetTokenInformation":               {},
	"GetUserNameW":                      {},
	"GetVersionExA":                     {},
	"GetVersionExW":                     {},
	"GetVolumeInformationA":             {},
	"GetVolumeInformationW":             {},
	"GetVolumeNameForVolumeMountPointW": {},
	"GetWindowLongPtrW":                 {},
	"GetWindowsDirectoryA":              {},
	"GlobalAddAtomW":                    {},
	"GlobalDeleteAtom":                  {},
	"GlobalGetAtomNameW":                {},
	"HeapAlloc":                         {},
	"HeapCompact":                       {},
	"HeapCreate":                        {},
	"HeapDestroy":                       {},
	"HeapFree":                          {},
	"HeapReAlloc":                       {},
	"HeapSize":                          {},
	"HeapValidate":                      {},
	"IN6_ADDR_EQUAL":                    {},
	"IN6_IS_ADDR_V4MAPPED":              {},
	"ImpersonateSelf":                   {},
	"InitializeAcl":                     {},
	"InitializeCriticalSection":         {},
	"InitializeSid":                     {},
	"IsDebuggerPresent":                 {},
	"IsWindow":                          {},
	"KillTimer":                         {},
	"LeaveCriticalSection":              {},
	"LoadLibraryA":                      {},
	"LoadLibraryExW":                    {},
	"LoadLibraryW":                      {},
	"LocalFree":                         {},
	"LockFile":                          {},
	"LockFileEx":                        {},
	"MapViewOfFile":                     {},
	"MessageBeep":                       {},
	"MessageBoxW":                       {},
	"MoveFileW":                         {},
	"MsgWaitForMultipleObjectsEx":       {},
	"MultiByteToWideChar":               {},
	"NetApiBufferFree":                  {},
	"NetGetDCName":                      {},
	"NetUserGetInfo":                    {},
	"OpenEventA":                        {},
	"OpenProcessToken":                  {},
	"OpenThreadToken":                   {},
	"OutputDebugStringA":                {},
	"OutputDebugStringW":                {},
	"PeekConsoleInputW":                 {},
	"PeekMessageW":                      {},
	"PeekNamedPipe":                     {},
	"PostMessageW":                      {},
	"PostQuitMessage":                   {},
	"PurgeComm":                         {},
	"QueryPerformanceCounter":           {},
	"QueryPerformanceFrequency":         {},
	"RaiseException":                    {},
	"ReadConsoleW":                      {},
	"ReadFile":                          {},
	"RegCloseKey":                       {},
	"RegConnectRegistryW":               {},
	"RegCreateKeyExW":                   {},
	"RegDeleteKeyW":                     {},
	"RegDeleteValueW":                   {},
	"RegEnumKeyExW":                     {},
	"RegEnumValueW":                     {},
	"RegOpenKeyExW":                     {},
	"RegQueryValueExW":                  {},
	"RegSetValueExW":                    {},
	"RegisterClassExW":                  {},
	"RegisterClassW":                    {},
	"RemoveDirectoryW":                  {},
	"ResetEvent":                        {},
	"RevertToSelf":                      {},
	"RtlGetVersion":                     {},
	"SearchPathW":                       {},
	"SendMessageTimeoutW":               {},
	"SendMessageW":                      {},
	"SetCommState":                      {},
	"SetCommTimeouts":                   {},
	"SetConsoleCtrlHandler":             {},
	"SetConsoleMode":                    {},
	"SetConsoleTextAttribute":           {},
	"SetCurrentDirectoryW":              {},
	"SetEndOfFile":                      {},
	"SetErrorMode":                      {},
	"SetEvent":                          {},
	"SetFileAttributesW":                {},
	"SetFilePointer":                    {},
	"SetFileTime":                       {},
	"SetHandleInformation":              {},
	"SetNamedSecurityInfoA":             {},
	"SetThreadPriority":                 {},
	"SetTimer":                          {},
	"SetWindowLongPtrW":                 {},
	"SetupComm":                         {},
	"Sleep":                             {},
	"SleepEx":                           {},
	"SystemTimeToFileTime":              {},
	"TerminateThread":                   {},
	"TranslateMessage":                  {},
	"UnlockFile":                        {},
	"UnlockFileEx":                      {},
	"UnmapViewOfFile":                   {},
	"UnregisterClassW":                  {},
	"WSAAsyncSelect":                    {},
	"WSAGetLastError":                   {},
	"WSAStartup":                        {},
	"WaitForInputIdle":                  {},
	"WaitForSingleObject":               {},
	"WaitForSingleObjectEx":             {},
	"WideCharToMultiByte":               {},
	"WriteConsoleW":                     {},
	"WriteFile":                         {},
	"WspiapiFreeAddrInfo":               {},
	"WspiapiGetAddrInfo":                {},
	"WspiapiGetNameInfo":                {},
	"_IO_putc":                          {},
	"_InterlockedCompareExchange":       {},
	"_InterlockedExchange":              {},
	"___errno_location":                 {},
	"__acrt_iob_func":                   {},
	"__assert_fail":                     {},
	"__atomic_load_n":                   {},
	"__atomic_store_n":                  {},
	"__builtin___memcpy_chk":            {},
	"__builtin___memmove_chk":           {},
	"__builtin___memset_chk":            {},
	"__builtin___snprintf_chk":          {},
	"__builtin___sprintf_chk":           {},
	"__builtin___strcat_chk":            {},
	"__builtin___strcpy_chk":            {},
	"__builtin___strncpy_chk":           {},
	"__builtin___vsnprintf_chk":         {},
	"__builtin_abort":                   {},
	"__builtin_abs":                     {},
	"__builtin_add_overflow":            {},
	"__builtin_add_overflowInt64":       {},
	"__builtin_add_overflowUint32":      {},
	"__builtin_add_overflowUint64":      {},
	"__builtin_bswap16":                 {},
	"__builtin_bswap32":                 {},
	"__builtin_bswap64":                 {},
	"__builtin_clzll":                   {},
	"__builtin_constant_p_impl":         {},
	"__builtin_copysign":                {},
	"__builtin_copysignf":               {},
	"__builtin_exit":                    {},
	"__builtin_expect":                  {},
	"__builtin_fabs":                    {},
	"__builtin_free":                    {},
	"__builtin_huge_val":                {},
	"__builtin_huge_valf":               {},
	"__builtin_inf":                     {},
	"__builtin_inff":                    {},
	"__builtin_isnan":                   {},
	"__builtin_malloc":                  {},
	"__builtin_memcmp":                  {},
	"__builtin_memcpy":                  {},
	"__builtin_memset":                  {},
	"__builtin_mmap":                    {},
	"__builtin_mul_overflow":            {},
	"__builtin_mul_overflowInt64":       {},
	"__builtin_mul_overflowUint128":     {},
	"__builtin_mul_overflowUint64":      {},
	"__builtin_nanf":                    {},
	"__builtin_object_size":             {},
	"__builtin_popcount":                {},
	"__builtin_prefetch":                {},
	"__builtin_printf":                  {},
	"__builtin_snprintf":                {},
	"__builtin_sprintf":                 {},
	"__builtin_strchr":                  {},
	"__builtin_strcmp":                  {},
	"__builtin_strcpy":                  {},
	"__builtin_strlen":                  {},
	"__builtin_sub_overflow":            {},
	"__builtin_sub_overflowInt64":       {},
	"__builtin_trap":                    {},
	"__builtin_unreachable":             {},
	"__ccgo_in6addr_anyp":               {},
	"__ccgo_sqlite3_log":                {},
	"__ctype_b_loc":                     {},
	"__ctype_get_mb_cur_max":            {},
	"__env_rm_add":                      {},
	"__errno_location":                  {},
	"__imp__environ":                    {},
	"__isalnum_l":                       {},
	"__isalpha_l":                       {},
	"__isdigit_l":                       {},
	"__islower_l":                       {},
	"__isnan":                           {},
	"__isnanf":                          {},
	"__isnanl":                          {},
	"__isoc99_sscanf":                   {},
	"__isprint_l":                       {},
	"__isspace_l":                       {},
	"__isxdigit_l":                      {},
	"__mingw_vfprintf":                  {},
	"__mingw_vfscanf":                   {},
	"__mingw_vfwprintf":                 {},
	"__mingw_vfwscanf":                  {},
	"__mingw_vsnprintf":                 {},
	"__mingw_vsnwprintf":                {},
	"__mingw_vsprintf":                  {},
	"__mingw_vsscanf":                   {},
	"__mingw_vswscanf":                  {},
	"__ms_vfscanf":                      {},
	"__ms_vfwscanf":                     {},
	"__ms_vscanf":                       {},
	"__ms_vsnprintf":                    {},
	"__ms_vsscanf":                      {},
	"__ms_vswscanf":                     {},
	"__ms_vwscanf":                      {},
	"__putenv":                          {},
	"__strchrnul":                       {},
	"_access":                           {},
	"_assert":                           {},
	"_beginthread":                      {},
	"_beginthreadex":                    {},
	"_byteswap_uint64":                  {},
	"_byteswap_ulong":                   {},
	"_chmod":                            {},
	"_chsize":                           {},
	"_commit":                           {},
	"_controlfp":                        {},
	"_endthreadex":                      {},
	"_errno":                            {},
	"_exit":                             {},
	"_fileno":                           {},
	"_findclose":                        {},
	"_findfirst32":                      {},
	"_findfirst64i32":                   {},
	"_findnext32":                       {},
	"_findnext64i32":                    {},
	"_fstat64":                          {},
	"_fstati64":                         {},
	"_ftime":                            {},
	"_imp___environ":                    {},
	"_isatty":                           {},
	"_localtime64":                      {},
	"_mkdir":                            {},
	"_msize":                            {},
	"_obstack_begin":                    {},
	"_obstack_newchunk":                 {},
	"_pclose":                           {},
	"_popen":                            {},
	"_putchar":                          {},
	"_set_abort_behavior":               {},
	"_setmode":                          {},
	"_snprintf":                         {},
	"_snwprintf":                        {},
	"_stat64":                           {},
	"_stati64":                          {},
	"_strdup":                           {},
	"_stricmp":                          {},
	"_strnicmp":                         {},
	"_unlink":                           {},
	"_vsnwprintf":                       {},
	"_wcsicmp":                          {},
	"_wcsnicmp":                         {},
	"_wopen":                            {},
	"_wunlink":                          {},
	"abort":                             {},
	"abs":                               {},
	"accept":                            {},
	"access":                            {},
	"acos":                              {},
	"acosh":                             {},
	"alarm":                             {},
	"asin":                              {},
	"asinh":                             {},
	"atan":                              {},
	"atan2":                             {},
	"atanh":                             {},
	"atexit":                            {},
	"atof":                              {},
	"atoi":                              {},
	"atol":                              {},
	"backtrace":                         {},
	"backtrace_symbols_fd":              {},
	"bind":                              {},
	"calloc":                            {},
	"ceil":                              {},
	"ceilf":                             {},
	"cfsetispeed":                       {},
	"cfsetospeed":                       {},
	"chdir":                             {},
	"chmod":                             {},
	"clock_gettime":                     {},
	"close":                             {},
	"closedir":                          {},
	"closesocket":                       {},
	"confstr":                           {},
	"connect":                           {},
	"copysign":                          {},
	"copysignf":                         {},
	"cos":                               {},
	"cosf":                              {},
	"cosh":                              {},
	"dlclose":                           {},
	"dlerror":                           {},
	"dlopen":                            {},
	"dlsym":                             {},
	"dup2":                              {},
	"environ":                           {},
	"execvp":                            {},
	"exit":                              {},
	"exp":                               {},
	"fabs":                              {},
	"fabsf":                             {},
	"fchmod":                            {},
	"fclose":                            {},
	"fcntl":                             {},
	"fcntl64":                           {},
	"fdopen":                            {},
	"ferror":                            {},
	"fflush":                            {},
	"fgetc":                             {},
	"fgets":                             {},
	"fileno":                            {},
	"floor":                             {},
	"fmod":                              {},
	"fopen":                             {},
	"fopen64":                           {},
	"fork":                              {},
	"fprintf":                           {},
	"fputc":                             {},
	"fputs":                             {},
	"fread":                             {},
	"free":                              {},
	"frexp":                             {},
	"fscanf":                            {},
	"fseek":                             {},
	"fstat":                             {},
	"fstat64":                           {},
	"fsync":                             {},
	"ftell":                             {},
	"ftruncate":                         {},
	"ftruncate64":                       {},
	"fts64_close":                       {},
	"fts64_open":                        {},
	"fts64_read":                        {},
	"fts_close":                         {},
	"fts_read":                          {},
	"fwrite":                            {},
	"gai_strerror":                      {},
	"gai_strerrorW":                     {},
	"getc":                              {},
	"getcwd":                            {},
	"getenv":                            {},
	"gethostname":                       {},
	"getpeername":                       {},
	"getpid":                            {},
	"getpwuid":                          {},
	"getrlimit":                         {},
	"getrlimit64":                       {},
	"getrusage":                         {},
	"getservbyname":                     {},
	"getsockname":                       {},
	"getsockopt":                        {},
	"gettimeofday":                      {},
	"gmtime":                            {},
	"gmtime_r":                          {},
	"htonl":                             {},
	"htons":                             {},
	"hypot":                             {},
	"inet_ntoa":                         {},
	"initstate_r":                       {},
	"ioctl":                             {},
	"ioctlsocket":                       {},
	"isalnum":                           {},
	"isalpha":                           {},
	"isatty":                            {},
	"isdigit":                           {},
	"islower":                           {},
	"isnan":                             {},
	"isnanf":                            {},
	"isnanl":                            {},
	"isprint":                           {},
	"isspace":                           {},
	"isxdigit":                          {},
	"kill":                              {},
	"ldexp":                             {},
	"link":                              {},
	"listen":                            {},
	"localtime":                         {},
	"localtime_r":                       {},
	"log":                               {},
	"log10":                             {},
	"lseek":                             {},
	"lseek64":                           {},
	"lstat":                             {},
	"lstat64":                           {},
	"lstrcmpiA":                         {},
	"lstrlenW":                          {},
	"malloc":                            {},
	"mblen":                             {},
	"mbstowcs":                          {},
	"mbtowc":                            {},
	"memchr":                            {},
	"memcmp":                            {},
	"memcpy":                            {},
	"memmove":                           {},
	"memset":                            {},
	"mkdir":                             {},
	"mkfifo":                            {},
	"mknod":                             {},
	"mkstemp64":                         {},
	"mkstemps":                          {},
	"mkstemps64":                        {},
	"mktime":                            {},
	"mmap":                              {},
	"mmap64":                            {},
	"modf":                              {},
	"mremap":                            {},
	"munmap":                            {},
	"ntohs":                             {},
	"obstack_free":                      {},
	"obstack_vprintf":                   {},
	"open":                              {},
	"open64":                            {},
	"opendir":                           {},
	"openpty":                           {},
	"pclose":                            {},
	"perror":                            {},
	"pipe":                              {},
	"popen":                             {},
	"posix_fadvise":                     {},
	"pow":                               {},
	"printf":                            {},
	"pselect":                           {},
	"putc":                              {},
	"putchar":                           {},
	"putenv":                            {},
	"puts":                              {},
	"qsort":                             {},
	"raise":                             {},
	"rand":                              {},
	"random":                            {},
	"random_r":                          {},
	"read":                              {},
	"readdir":                           {},
	"readlink":                          {},
	"readv":                             {},
	"realloc":                           {},
	"realpath":                          {},
	"recv":                              {},
	"rename":                            {},
	"rewind":                            {},
	"rmdir":                             {},
	"round":                             {},
	"select":                            {},
	"send":                              {},
	"setbuf":                            {},
	"setenv":                            {},
	"setlocale":                         {},
	"setmode":                           {},
	"setrlimit":                         {},
	"setrlimit64":                       {},
	"setsid":                            {},
	"setsockopt":                        {},
	"setvbuf":                           {},
	"shutdown":                          {},
	"sigaction":                         {},
	"sin":                               {},
	"sinf":                              {},
	"sinh":                              {},
	"sleep":                             {},
	"snprintf":                          {},
	"socket":                            {},
	"sprintf":                           {},
	"sqrt":                              {},
	"sscanf":                            {},
	"stat":                              {},
	"stat64":                            {},
	"stderr":                            {},
	"stdin":                             {},
	"stdout":                            {},
	"strcasecmp":                        {},
	"strcat":                            {},
	"strchr":                            {},
	"strcmp":                            {},
	"strcpy":                            {},
	"strcspn":                           {},
	"strdup":                            {},
	"strerror":                          {},
	"strlen":                            {},
	"strncmp":                           {},
	"strncpy":                           {},
	"strpbrk":                           {},
	"strrchr":                           {},
	"strstr":                            {},
	"strtod":                            {},
	"strtol":                            {},
	"strtoul":                           {},
	"symlink":                           {},
	"sysconf":                           {},
	"system":                            {},
	"tan":                               {},
	"tanh":                              {},
	"tcgetattr":                         {},
	"tcsendbreak":                       {},
	"tcsetattr":                         {},
	"time":                              {},
	"timezone":                          {},
	"tolower":                           {},
	"toupper":                           {},
	"trunc":                             {},
	"tzset":                             {},
	"umask":                             {},
	"uname":                             {},
	"ungetc":                            {},
	"unlink":                            {},
	"unsetenv":                          {},
	"usleep":                            {},
	"utime":                             {},
	"utimes":                            {},
	"vasprintf":                         {},
	"vfprintf":                          {},
	"vprintf":                           {},
	"vsnprintf":                         {},
	"vsprintf":                          {},
	"waitpid":                           {},
	"wcrtomb":                           {},
	"wcschr":                            {},
	"wcscmp":                            {},
	"wcscpy":                            {},
	"wcsicmp":                           {},
	"wcslen":                            {},
	"wcsncmp":                           {},
	"wcsrtombs":                         {},
	"wcstombs":                          {},
	"wctomb":                            {},
	"wcwidth":                           {},
	"write":                             {},
	"wsprintfA":                         {},
	"wsprintfW":                         {},
}
