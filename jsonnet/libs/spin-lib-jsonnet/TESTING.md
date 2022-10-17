# Testing

```shell
jt
```

Test-file Example:

```jsonnet
// Import the libsonnet to test.
local myLib = import '../myLib.libsonnet';

// Add what will be rendered to the "tests" array
local tests = [
    myLib.doSomething('potato'),
    myLib.doSomethingElse('bread'),
    myLib.SomeObject {fruit: 'apple'},
];

// Add how the tests are expected to be rendered to the "results" array
local assertions = [
    ['fries', 'baked potato', 'mushed potatoes'],
    'sandwhich',
    {fruit: 'apple', price: '$1.00'}
];

// Render/Provide back this object - checking if the tests match the expected results, and providing their values back.
{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
```
