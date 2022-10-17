/**
    @file Contains objects  to create Spinnaker pipeline Parameters.
**/

/**
    Creates a Parameter that can be used by the Spinnaker Pipeline.

    @example
        parameter.Parameter {
            name: 'my-param'
        }

        parameter.Parameter {
            name: 'my-param',
            description: 'Does magic',
            label: 'My Parameter'
        }

        parameter.Parameter {
            name: 'my-param',
            hasOptions: true,
            options = [
                {value: "potato"},
                {value: "apple"},
            ]
        }

    @property {String} name - Name of the parameter.
    @property {String} [default=''] - Default value of the parameter.
    @property {String} [description=''] - Description for the parameter.
    @property {Boolean} [hasOptions=false] - If the parameter has option/is a drop-down choice box.
    @property {Array<Object>} [options=[]] - Options for the parameter. Toggled by `hasOptions` property.
    @property {String} [label=''] - The human-friendly name/label for this parameter.
    @property {String} [pinned=false] - Whether or not this parameter is pinned.
    @property {String} [required=false] - Whether or not this parameter is required.
**/
local Parameter = {
  name: error '`name` (String) property is required for Parameter',

  default: '',
  description: '',
  hasOptions: false,
  label: '',
  options: [],
  pinned: false,
  required: false,

};

// Exposed for public use.
{
  Parameter:: Parameter,
}
