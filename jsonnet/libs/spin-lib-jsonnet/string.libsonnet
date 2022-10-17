/**
    @file Contains methods for operating on Strings.
    @todo Revisit the module in a new PR to check what is necessary for `spin-lib-jsonnet` to support.
**/

/**
    Trims a given string, removing any blank spaces or new lines at the start or end of the string.

    @param {String} str - String that should be trimmed.

    @returns {String} A trimmed string, without any leading or trailing blank spaces or new lines.
**/
local trim(str) = (
  if std.startsWith(str, ' ') || std.startsWith(str, '\n') then
    trim(std.substr(str, 1, std.length(str)))
  else if std.endsWith(str, ' ') || std.endsWith(str, '\n') then
    trim(std.substr(str, 0, std.length(str) - 1))
  else
    str
);

/**
    Capitalize a given string.

    @todo Figure out if this belongs in the standard lib.

    @example
        string.capitalize("potato")     // returns "Potato"

    @param {String} stringToCapitalize - String to capitalize.

    @returns {String} The same string, but capitalized.
**/
local capitalize(stringToCapitalize) = (
  local firstChar = std.asciiUpper(std.stringChars(stringToCapitalize)[0]);
  firstChar + std.substr(stringToCapitalize, 1, std.length(stringToCapitalize))
);

/**
    Returns a string of an array in a "pretty format".

    @example
        string.arrayToPrettyString([apples, oranges, tomatoes, potatoes], ', ')   // returns "apples, oranges, tomatoes, potatoes"
        string.arrayToPrettyString([apples, oranges, tomatoes, potatoes], ', ')   // returns "apples, oranges, tomatoes, potatoes"
        string.arrayToPrettyString([apples, oranges, tomatoes, potatoes], '-')   // returns "apples-oranges-tomatoes-potatoes"

    @param {Array<String>} givenArray - The given array to make a pretty string of.
    @param {String} seperator - The seperator used between items. Default `,`.

    @returns {String} The string representation of a given array, in a pretty format.
**/
local arrayToPrettyString(givenArray, seperator=',') = (
  local arrayString =
    std.strReplace(
      std.strReplace(
        std.strReplace(
          std.manifestJsonEx(givenArray, ''),
          '\n',
          ''
        ),
        '"',
        ''
      ),
      ',',
      seperator
    );

  if givenArray != null && std.length(givenArray) > 0 then
    trim(std.substr(arrayString, 1, std.length(arrayString) - 2))
  else ''
);

// Exposed for public use.
{
  arrayToPrettyString:: arrayToPrettyString,
  capitalize:: capitalize,
  trim:: trim,
}
