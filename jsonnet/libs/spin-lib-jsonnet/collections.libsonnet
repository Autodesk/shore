/**
    @file Contains methods for operating on collections.
**/

/**
    Find the correct key by recursing through the object's keys validating they exists and returning the value at the end.

    If no value was found, return the default.


    For "key" in "keys", test that obj["key"] exists.

    If obj["key"] exists, pop the first element by creating a new list without it.

    Recurse through the Object using Keys until either a match is found or "Default" is returned.

    @param {Object} obj - Object to get the key from.
    @param {Array<String>} keysPath - An array of keys.
    @param {Object} defValue - The default value to provide.

    @returns {Object} The object from the map, or the default.
**/
local getNestedValueSafely(obj, keysPath, defValue) = (
  assert std.isObject(obj) : 'The `obj` argument in `getNestedValueSafely` must be of type `Object`';
  assert std.isArray(keysPath) : 'The `keys` argument in `getNestedValueSafely` must be of type `Array`';
  local currentKey = keysPath[0];

  if std.objectHas(obj, currentKey) then
    // Use Array slice to get a subset of
    local recursionKeys = keysPath[1:std.length(keysPath)];

    if std.length(recursionKeys) == 0 then
      obj[currentKey]
    else
      getNestedValueSafely(obj[currentKey], recursionKeys, defValue)
  else
    defValue
);

/**
    Attempts to get an object using a key from given object.

    If the given object does not have a key, it will retusrn the supplied default.

    Does look at hidden/private (::) fields.

    @param {Object} givenMap - The given object that will be accessed with the given key.
    @param {String} key - The given key for the value/object to get.
    @param {Object} defValue - The default to return

    @returns {Object} The object/value for the given key or the supplied default
**/
local getObjectAllSafely(givenMap, key, defValue) = (
  if std.objectHasAll(givenMap, key) then
    givenMap[key]
  else
    defValue
);

/**
    Checks if the given object has the given key - which is also not empty ('') nor null.

    @param {Object} givenMap - The given object that will be checked.
    @param {String} key - The given key to check.

    @returns {Boolean} True if the object has the key that is not empty and not null, false if otherwise.
**/
local isNotNullOrEmpty(givenMap, key) = (
  if std.objectHas(givenMap, key) then
    if givenMap[key] == '' || givenMap[key] == null then
      false
    else
      true
  else
    false
);

/**
    Checks if the given object or string is null or empty

    @param {any} Item to test
    @param {boolean} true: item is empty/null, false: item is not empty/null 

    @returns {Boolean} True if the object is null or empty
**/
local isNullOrEmpty(val) = (
    val == null || val == {} || val == '' || (std.isArray(val) && std.length(val) == 0)
);

// Exposed for public use.
{
  getNestedValueSafely:: getNestedValueSafely,
  getObjectAllSafely:: getObjectAllSafely,
  isNotNullOrEmpty:: isNotNullOrEmpty,
  isNullOrEmpty:: isNullOrEmpty,
}
