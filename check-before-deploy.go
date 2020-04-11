package main

import (
	"flag"
	"fmt"
	"strings"
	"os"
	"io/ioutil"
    "path/filepath"
	"gopkg.in/yaml.v3"
//	"code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[1;35m%s\033[0m"
	OkColor   = "\033[1;32m%s\033[0m"
)

const c_TypeExisting = "org.cloudfoundry.existing-service"
const c_TypeManaged = "org.cloudfoundry.managed-service"
const c_CostPaid = "paid"

type YalmData struct {
    ID string `yaml:"ID"`
	Resources []struct{
		Name       string `yaml:"name"`
        Type       string `yaml:"type"`
		Parameters struct{
			ServicePlan    string `yaml:"service-plan"`
			Service	       string `yaml:"service"`
		}`yaml:"parameters"`
	}`yaml:"resources"`
}


type BindingResults []struct {
	Binding string
	Status bool
}


type PluginParams struct {
	file *bool
	yalmData YalmData
}

func main() {
	plugin.Start(new(PluginParams))
}

func (pluginDemo *PluginParams) Run(cliConnection plugin.CliConnection, args []string) {
	// Initialize flags

	// GetService(serviceInstance string) (plugin_models.GetService_Model, error)
	
	pluginFlag := flag.NewFlagSet("check-before-deploy", flag.ExitOnError)
	file := pluginFlag.String("file", "f", "--file path/to/some/file.yalm")
	checkbinding := pluginFlag.Bool("check-binding", false, "Check the YALM file if the binded services exist in org / space")
	checkservice := pluginFlag.Bool("check-service", false, "Check the YALM file if exist services plant exist in org/space")
	allChecks := pluginFlag.Bool("all", false, "Active all validations")


	// Parse starting from [1] because the [0]th element is the
	// Check parameter file
	err := pluginFlag.Parse(args[1:])
	if err != nil || file == nil || *file == "" || *file == "f" {
		fmt.Println("Error parameter file obligatory - please use -file /your_path/file.yalm")
		os.Exit(1)
	}
	
	if file != nil {
		//Get and parse file
		pluginDemo.ReadFile(*file)
		//check binding
		if *checkbinding || *allChecks {
			pluginDemo.yalmData.CheckResourceListBinding(cliConnection)
		}
		if *checkservice || *allChecks {
			pluginDemo.yalmData.CheckResourceListPlans(cliConnection)
		}
		//fmt.Printf("Value: %#v\n", pluginDemo.yalmData.Resources)
		fmt.Println("")
		fmt.Println("")
	}

	
}

// Read Yalm file
func (pluginDemo *PluginParams) ReadFile(file string) {

	// get file
    filename, _ := filepath.Abs(file)
    yamlFile, err := ioutil.ReadFile(filename)
	
	if err != nil {
		fmt.Printf(ErrorColor,err)
		os.Exit(1)
    }
	
	// parse file
	err = yaml.Unmarshal(yamlFile, &pluginDemo.yalmData)
    if err != nil {
		fmt.Printf(ErrorColor,err)
		os.Exit(1)
    }
	
}

// Check bindings
func (yalmParsed YalmData) CheckResourceListPlans(cliConnection plugin.CliConnection) {

	//var services plugin_models.GetService_Model
	//var errorCliCommand error
	var errorServices bool 
	var errorServicesPlan bool 

	fmt.Println("")
	fmt.Printf(InfoColor, "  Check plans")
	fmt.Println("")
	fmt.Printf(InfoColor, "-------------------------------------------------------------")
	fmt.Println("")

    for _, resource := range yalmParsed.Resources {
		//Check only 
		if resource.Type ==  c_TypeManaged{
			command_result, errorCliCommand := cliConnection.CliCommandWithoutTerminalOutput("marketplace", "-s", resource.Parameters.Service)
			if errorCliCommand == nil {
				errorServicesPlan = true
				for i := 4; i < len(command_result); i++ {
					//fmt.Println(command_result[i])
					words := strings.Fields(command_result[i])
					if resource.Parameters.ServicePlan == words[0]{
						fmt.Printf(OkColor,resource.Name)
						if words[len(words)-1] == c_CostPaid{
							fmt.Printf(WarningColor," - [Paid]")
						}
						fmt.Println("")
						errorServicesPlan = false
						break
						//fmt.Println(words[0], words[len(words)-1])
					}
				}
				if errorServicesPlan{
					fmt.Printf(ErrorColor,resource.Name)
					fmt.Printf(" - ")
					fmt.Printf(resource.Parameters.ServicePlan)
					fmt.Printf(" service plan does not exist");
					fmt.Println("")
				}
			}else{
				fmt.Printf(ErrorColor,resource.Name )
				fmt.Printf(" - ")
				fmt.Printf(resource.Parameters.Service)
				fmt.Printf(" service does not exist")
				fmt.Println("")
				errorServices = true
			}

		}
	}
	
	if  errorServices == false {
		fmt.Println(" > Services correct")
	}

}

// Check bindings
func (yalmParsed YalmData) CheckResourceListBinding(cliConnection plugin.CliConnection) {

	/*
	fmt.Printf(InfoColor, "Info")
	fmt.Println("")
	fmt.Printf(NoticeColor, "Notice")
	fmt.Println("")
	fmt.Printf(WarningColor, "Warning")
	fmt.Println("")
	fmt.Printf(ErrorColor, "Error")
	fmt.Println("")
	fmt.Printf(DebugColor, "Debug")
	fmt.Println("")
	fmt.Printf(OkColor, "OkColor")
	fmt.Println("")
	*/

	//var services plugin_models.GetService_Model
	var errorCliCommand error
	var errorBinding bool 

	fmt.Println("")
	fmt.Printf(InfoColor, "  Check binding")
	fmt.Println("")
	fmt.Printf(InfoColor, "-------------------------------------------------------------")
	fmt.Println("")

    for _, resource := range yalmParsed.Resources {
		//Check only 
		if resource.Type ==  c_TypeExisting{
			//services,errorCliCommand = cliConnection.GetService(resource.Name)
			_,errorCliCommand = cliConnection.GetService(resource.Name)
			if errorCliCommand != nil {
				errorBinding = true
				fmt.Printf(ErrorColor,resource.Name)
				fmt.Printf(" [ service does not exist in your Org/Space] ")
				fmt.Println("")
			}else{
				fmt.Printf(OkColor,resource.Name)
				fmt.Println("")
			}
			//fmt.Println(resource.Name)
		}
	}
	
	if  errorBinding == false {
		fmt.Println(" > Bindings correct")
	}

}

func (pluginDemo *PluginParams) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "check-before-deploy",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 6,
		},
		Commands: []plugin.Command{
			{
				Name:     "check-before-deploy",
				Alias:    "cbd",
				HelpText: "Check the Yalm file and services before deploying your MTA",
				UsageDetails: plugin.Usage{
					Usage: "cf check-before-deploy -file [path] -check-binding -check-service -all",
					Options: map[string]string{
						"file": "obligatory - Path with YALM file",
						"check-binding": "Check the YALM file if the binded services exist in org / space",
						"check-service": "Check the YALM file if exist services plant exist in org/space",
						"all": "Active all validations",
					},
				},
			},
		},
	}
}