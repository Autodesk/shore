local string = import '../string.libsonnet';

// Data used for tests
local multiLineText = |||

  potato

|||;

local tests = [
  string.trim(multiLineText),
  string.capitalize('potato'),

  string.arrayToPrettyString(['potato']),
  string.arrayToPrettyString(['potato', 'tomato', 'apple', 'orange']),
  string.arrayToPrettyString(['potato', 'tomato', 'apple', 'orange'], ', '),
  string.arrayToPrettyString(['potato', 'tomato', 'apple', 'orange'], '-'),
  string.arrayToPrettyString([]),
  string.arrayToPrettyString(''),
  string.arrayToPrettyString(null),
];

local assertions = [
  'potato',
  'Potato',

  'potato',
  'potato,tomato,apple,orange',
  'potato, tomato, apple, orange',
  'potato-tomato-apple-orange',
  '',
  '',
  '',
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
