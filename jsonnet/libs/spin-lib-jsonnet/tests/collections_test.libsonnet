local collections = import '../collections.libsonnet';

// Data used for tests.
local nestedObject = {
  test1: {
    test2: {
      test3: {
        test4: 'value',
      },
    },
  },
};

local complexNestedObject = {
  'test-l1': {
    'test-l2-1': [
      1,
      2,
      3,
      4,
      5,
      5,
      6,
      67,
      7,
      78,
      6,
      5,
      43,
    ],
    'test-l2-2': {
      Underalted: 'Code',
    },
    'test-l2-3': {
      These: {
        Are: {
          Not: {
            The: {
              Droids: {
                You: {
                  Are: {
                    Looking: {
                      For: false,
                    },
                  },
                },
              },
            },
          },
        },
      },
    },
  },
};

local testDefaultValue = {
  So: {
    Close: {
      But: {
        Banana: true,
      },
    },
  },
};

local tests = [
  collections.getObjectAllSafely({ apple:: { color: 'red', taste: 'sweet' } }, 'apple', 'Not found'),
  collections.getObjectAllSafely({ apple:: { color: 'red', taste: 'sweet' } }, 'orange', 'Not found'),

  collections.isNotNullOrEmpty({ apple: { color: 'red', taste: 'sweet' } }, 'apple'),
  collections.isNotNullOrEmpty({ apple: { color: 'red', taste: 'sweet' } }, 'orange'),
  collections.isNotNullOrEmpty({ orange: '' }, 'orange'),
  collections.isNotNullOrEmpty({ orange: null }, 'orange'),

  collections.getNestedValueSafely(nestedObject, ['test1', 'test2', 'test3', 'test4'], 'defaultValue'),
  collections.getNestedValueSafely(complexNestedObject, ['test-l1', 'test-l2-3', 'These', 'Are', 'Not', 'The', 'Droids', 'You', 'Are', 'Looking', 'For'], 'defaultValue'),
  collections.getNestedValueSafely(testDefaultValue, ['So', 'Close', 'But', 'No', 'Apples?'], 'Should Have Been Banana!!'),

  collections.isNullOrEmpty({}),
  collections.isNullOrEmpty([]),
  collections.isNullOrEmpty([1]),
];

local assertions = [
  { color: 'red', taste: 'sweet' },
  'Not found',

  true,
  false,
  false,
  false,

  'value',
  false,
  'Should Have Been Banana!!',

  true,
  true,
  false
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
