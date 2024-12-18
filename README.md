![greenfra-white-bg](https://github.com/user-attachments/assets/6012774f-d0ae-4cbf-860a-31c9875bcdbb)

# Greenfra - Greening your infrastructure one line at a time

*Powered with Go* 🐹

Estimate the environmental impact of your infrastructure changes within your Terraform projects!

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Commands](#commands)
  - [Flags](#flags)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Introduction
Greenfra is a tool designed to help you estimate the environmental impact of the changes you make to your infrastructure using Terraform. By analyzing your local Terraform configurations, Greenfra provides insights into the potential carbon footprint of your planned resource changes.

Greenfra is currently under development.

## Features
- Estimate carbon footprint for various Terraform resources.
- Analyze local Terraform changes before applying them.
- Support for various instance types and resources.
- Easy to use command-line interface.

## Installation

Use Homebrew: 

```sh 
brew tap donovan1905/greenfra 
brew install Donovan1905/homebrew-greenfra/greenfra
```

To install Greenfra, you need to have Go installed on your system. Then, you can clone the repository and build the project:

```sh
git clone https://github.com/Donovan1905/greenfra.git
cd greenfra
go build
```

Or use the docker image provided here : https://hub.docker.com/r/donovan190500/greenfra

## Usage

### Commands
Greenfra supports the following commands:

- `analyze` : Estimate the carbon footprint for EC2 changes in your terraform project.

### Flags
Greenfra supports the following flags:

- `-exec-plan` : Choose to whether or not asking greenfra to execute the terraform plan. If not used, you will need to provide a path to a terraform plan file. You can also set the environnement variable GREENFRA_EXEC_PLAN to true to alway execute terraform plan with Greenfra.

### Comments 

You can add greenfra comment to your terraform files to provide additional information about your resources (number of execution per month, mean execution duration, etc...).

**Currently, this comment are supported :** 

Lambda : 
- monthly_invocation : the number of Lambda invocations expected per month for this function
- mean_execution_time : the mean execution time expected for this function in milliseconds

EC2 : 
- usage_percentage : the percentage of an EC2 instance usage per month (50 % = 15 days)

The format is the following : 

```terraform
/* greenfra
<key>=<value>
*/
```

Example for lambda : 
```terraform
/* greenfra
monthly_invocation=1000000
mean_execution_time=300
*/
```

### Example Commands
Here are some example commands to use Greenfra:

- To estimate the carbon footprint of your instances without using an existing plan:
  ```sh
  greenfra -exec-plan analyze
  ```

- To estimate the carbon footprint of your instances with an already generated plan:
  ```sh
  greenfra analyze <relative path to your plan file>
  ```

## Contributing
We welcome contributions from the community! Here’s how you can help:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -m 'Add some feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a new Pull Request.


## License
Greenfra is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for more information.
