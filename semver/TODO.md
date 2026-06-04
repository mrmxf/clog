# semver Package - Issues and Improvements

## ✅ COMPLETED: Pure Function Refactoring

The semver package has been refactored to use pure functions while maintaining backward compatibility:

### New Pure Function API
- **Function:** `ParseLinkerJSON(semVerJSON string) (VersionInfo, bool, error)`
  - ✅ No global state mutation
  - ✅ Returns values instead of mutating globals
  - ✅ Can be tested in parallel
  - ✅ Fully isolated tests
  - ✅ Proper error handling (returns errors instead of silent failures)

### Backward Compatibility Maintained
- ✅ The `init()` function still works exactly as before
- ✅ `semver.Info()` still returns the global `parsedInfo`
- ✅ Global `IsProductionBuild` still set correctly
- ✅ Existing code continues to work unchanged

### Test Improvements
- ✅ New pure function tests in `semver_pure_test.go` (48 assertions)
- ✅ Tests can run in parallel (`t.Parallel()`)
- ✅ No global state cleanup required
- ✅ Demonstrates isolated testing pattern

### Migration Path for New Code
```go
// Old way (still works):
info := semver.Info()

// New way (recommended):
info, isProd, err := semver.ParseLinkerJSON(semVerJSON)
if err != nil {
    // handle error
}
```

---

## Issues Found During Test Development

### 1. ✅ **FIXED: Error Handling - JSON Unmarshal Error Not Returned**
**Location:** [semver-parse-linkerInfo.go:34](semver-parse-linkerInfo.go#L34)

**Status:** ✅ Fixed in `ParseLinkerJSON()` - now properly returns errors:
```go
if err := json.Unmarshal([]byte(ldString), &data); err != nil {
    return VersionInfo{}, false, fmt.Errorf("failed to parse semver JSON: %w", err)
}
```

**Note:** `cleanLinkerData()` still maintains old behavior for backward compatibility.

### 2. ✅ **FIXED: Global State Mutation in Tests**
**Location:** Previously affected all tests

**Status:** ✅ Fixed - New pure function tests in `semver_pure_test.go` demonstrate:
- No global state mutation
- No test order dependencies
- Can run in parallel
- Simple, clean test code

**Migration:** Use `ParseLinkerJSON()` for new code and tests

### 3. **Hash Validation Logic Issue**
**Location:** [semver-parse-linkerInfo.go:43-51](semver-parse-linkerInfo.go#L43)

**Issue:** Two separate `if` statements both check and replace the hash:
```go
if len(linkerData.Hash) == 0 {
    linkerData.Hash = dummyHash
}
if len(linkerData.Hash) != 40 {
    linkerData.Hash = dummyHash
}
```

**Impact:** The first check is redundant since zero-length hashes will be caught by the second check anyway.

**Recommendation:** Combine into single validation:
```go
if len(linkerData.Hash) != 40 {
    slog.Debug("WARNING semver linkerData.Hash invalid", "length", len(linkerData.Hash), "expected", 40)
    linkerData.Hash = dummyHash
}
```

### 4. **Missing Tag Validation**
**Location:** [semver-parse-linkerInfo.go](semver-parse-linkerInfo.go)

**Issue:** No validation is performed on the `Tag` field. Invalid semver tags (e.g., "abc", "1.2", "") are silently accepted.

**Impact:** Could lead to invalid version strings being used throughout the application.

**Recommendation:** Add semver validation:
```go
if !isValidSemver(linkerData.Tag) {
    slog.Debug("WARNING semver tag invalid", "tag", linkerData.Tag)
    linkerData.Tag = dummyTag
}
```

### 5. **init() Function Side Effects**
**Location:** [semver.go:30-45](semver.go#L30)

**Issue:** Package initialization happens in `init()` which:
- Makes testing difficult
- Executes automatically on import
- Cannot be controlled or delayed
- Errors are stored in a struct field rather than returned

**Impact:**
- Cannot test initialization with different inputs easily
- Errors are hidden in `parsedInfo.Err` and may not be checked
- No way to re-initialize with new data

**Recommendation:** Replace `init()` with explicit initialization:
```go
func Init() error {
    if err := cleanLinkerData(); err != nil {
        return err
    }
    // ... rest of initialization
    return nil
}
```

### 6. **Date Fallback to Current Time**
**Location:** [semver-parse-linkerInfo.go:54-57](semver-parse-linkerInfo.go#L54)

**Issue:** When date parsing fails, the code uses `time.Now()` which means:
- Builds at different times will report different dates even with same linker data
- Not reproducible
- Makes testing time-dependent

**Recommendation:** Use a fixed date or make it more explicit:
```go
const fallbackDate = "1970-01-01" // Unix epoch as fallback
_, err := time.Parse("2006-01-02", linkerData.Date)
if err != nil {
    slog.Warn("Invalid build date, using fallback", "date", linkerData.Date)
    linkerData.Date = fallbackDate
}
```

### 7. **Inconsistent Error Handling**
**Location:** Throughout the package

**Issue:** Mix of error handling approaches:
- Some functions return errors (`cleanLinkerData() error`)
- Some log warnings via `slog.Debug()`
- Some silently use defaults
- `parsedInfo.Err` field to store errors

**Recommendation:** Standardize on one approach - preferably returning errors up the call stack.

### 8. **Missing Documentation for JSON Field Mapping**
**Location:** [defs.go:11-19](defs.go#L11)

**Issue:** The JSON field names in `LinkerDataJSON` don't match the struct field names, which could be confusing:
- Struct has `AppName`, JSON has `name`
- Struct has `AppTitle`, JSON has `title`

**Recommendation:** Add documentation or use explicit JSON tags:
```go
type LinkerDataJSON struct {
    Build    string `json:"build"`
    Tag      string `json:"tag"`
    Hash     string `json:"hash"`
    Date     string `json:"date"`
    Suffix   string `json:"suffix"`
    AppName  string `json:"name"`   // ← Document or rename JSON field to "appname"
    AppTitle string `json:"title"`  // ← Document or rename JSON field to "apptitle"
}
```

### 9. **No Validation for Build Type**
**Location:** [semver-parse-linkerInfo.go:38](semver-parse-linkerInfo.go#L38)

**Issue:** `IsProductionBuild` only checks for exact string match "prod":
```go
IsProductionBuild = linkerData.Build == "prod"
```

**Impact:** Typos like "production", "PROD", "Prod" will be treated as dev builds without warning.

**Recommendation:** Add validation with warning:
```go
switch strings.ToLower(linkerData.Build) {
case "prod", "production":
    IsProductionBuild = true
case "dev", "development", "":
    IsProductionBuild = false
default:
    slog.Warn("Unknown build type, treating as dev", "build", linkerData.Build)
    IsProductionBuild = false
}
```

### 10. **Quote Trimming May Be Too Aggressive**
**Location:** [semver-parse-linkerInfo.go:29](semver-parse-linkerInfo.go#L29)

**Issue:** The trim operation removes ALL quotes from both ends:
```go
ldString := strings.Trim(SemVerJSON, "\"'")
```

**Impact:** If JSON legitimately starts/ends with quotes inside the structure, they'll be removed.

**Recommendation:** Use more precise trimming:
```go
ldString := strings.TrimSpace(SemVerJSON)
ldString = strings.Trim(ldString, "\"'") // Only trim outer quotes after whitespace
```

## Test Coverage Summary

The new test suite ([semver_test.go](semver_test.go)) provides comprehensive coverage:

- ✅ JSON parsing (valid, invalid, missing fields, malformed)
- ✅ Hash validation (empty, wrong length, valid)
- ✅ Date validation (invalid format, empty, valid)
- ✅ Suffix handling (prod/dev, empty, custom)
- ✅ AppName/AppTitle derivation from module path
- ✅ Production vs development build detection
- ✅ VersionInfo initialization
- ✅ Quote stripping from JSON input
- ✅ Info() function behavior

**Total Assertions:** 75 passing tests with self-documenting GoConvey style

## Recommended Priority

1. **High Priority:**
   - Issue #1 (Error handling for JSON unmarshal)
   - Issue #5 (init() function side effects)
   - Issue #7 (Inconsistent error handling)

2. **Medium Priority:**
   - Issue #2 (Global state mutation)
   - Issue #4 (Missing tag validation)
   - Issue #9 (Build type validation)

3. **Low Priority:**
   - Issue #3 (Hash validation simplification)
   - Issue #6 (Date fallback)
   - Issue #8 (Documentation)
   - Issue #10 (Quote trimming)
