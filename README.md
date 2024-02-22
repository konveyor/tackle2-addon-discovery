# tackle2-addon-alizer

Automatically tags an application with language, framework, and tooling discovered by [devfiles/alizer](github.com/devfiles/alizer).

## Usage

1. Create Addon CR pointing to the tackle2-addon-alizer image. (See `hack/alizer_cr.yaml`)
2. Add application with source repository details to the Konveyor inventory.
3. Create Task pointing at the `alizer` addon and desired application. (See `hack/discovery.sh`)