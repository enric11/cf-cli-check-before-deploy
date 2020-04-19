package main

import (
	"flag"
	"fmt"
	"strings"
	"os"
	"io/ioutil"
    "path/filepath"
	"gopkg.in/yaml.v3"
	"code.cloudfoundry.org/cli/plugin"
	"archive/zip"
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

type YAMLData struct {
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
	YAMLData YAMLData
}

var resultsChecks struct {
	errorBindings bool
	errorServices bool
	errorServicesPlan bool
}

func main() {
	plugin.Start(new(PluginParams))
}

func (pluginDemo *PluginParams) Run(cliConnection plugin.CliConnection, args []string) {
	// Initialize flags
	pluginFlag := flag.NewFlagSet("check-before-deploy", flag.ExitOnError)
	file := pluginFlag.String("file", "f", "-file path/to/some/file.YAML")
	mta := pluginFlag.String("mta", "m", "-mta path/to/some/mta_file.mtar")
	checkbinding := pluginFlag.Bool("check-binding", false, "-check-binding > Check the YAML file if the binded services exist in org / space")
	checkservice := pluginFlag.Bool("check-service", false, "-check-service > Check the YAML file if exist services plant exist in org / space")
	allChecks := pluginFlag.Bool("all", false, "Active all validations")

	// Parse starting from [1] because the [0]th element is the
	// Check parameter file
	err := pluginFlag.Parse(args[1:])
	
	// Control to check errors when unisntall plugin in local
	if(args[0] == "CLI-MESSAGE-UNINSTALL"){
		os.Exit(0)
	}

	//check yaml file
	if *mta == "m" && *file != "f"   {
		//Read simple file
		pluginDemo.ReadFile(*file)
	}else if *file == "f" && *mta != "m"{
		//Read check mta file
		pluginDemo.UnzipMTA(*mta)	
	}else{
		fmt.Println("Parameter missing, please use: '-file /your_path/file' or '-mta /your_path/file'",err)
		os.Exit(1)
	}
	
	//check binding
	if *checkbinding || *allChecks {
		pluginDemo.YAMLData.CheckResourceListBinding(cliConnection)
	}
	if *checkservice || *allChecks {
		pluginDemo.YAMLData.CheckResourceListPlans(cliConnection)
	}

	if resultsChecks.errorBindings == true || resultsChecks.errorServices == true || resultsChecks.errorServicesPlan == true {
		fmt.Println("")
		fmt.Println("")
		fmt.Printf(ErrorColor,"Some errors have been detected, please fix them before deploy")
		fmt.Println("")
		os.Exit(1)
	}
}

// -----------------------------------------------------------------------------------------
func (pluginDemo *PluginParams)UnzipMTA(file string) {

	// Open a zip archive for reading.
	r, err := zip.OpenReader(file)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Close()

	// Iterate through the files in the archive,
	for _, f := range r.File {
			if f.Name == "META-INF/mtad.yaml" {
				rc, err := f.Open()
				if rc == nil {
					fmt.Println(err)
				}

				b1 := make([]byte,100000000)
				n1, err := rc.Read(b1)
				
				b2 := b1[:n1]
				rc.Close()
				pluginDemo.ParseYAMLFile(b2)
			}
	}
}
//--------------------------------------------------------------------------------------------

// Read YAML file
func (pluginDemo *PluginParams) ReadFile(file string) {

	// get file
    filename, _ := filepath.Abs(file)
    yamlFile, err := ioutil.ReadFile(filename)
	
	if err != nil {
		fmt.Printf(ErrorColor,err)
		os.Exit(1)
	}
	
	// parse file
	pluginDemo.ParseYAMLFile(yamlFile)		
}


func (pluginDemo *PluginParams) ParseYAMLFile(yamlFile []byte) {
		// parse file
		err := yaml.Unmarshal(yamlFile, &pluginDemo.YAMLData)
		if err != nil {
			fmt.Printf(ErrorColor,err)
			os.Exit(1)
		}
}

// Check bindings
func (YAMLParsed YAMLData) CheckResourceListPlans(cliConnection plugin.CliConnection) {

	var errorServicesPlan bool 

	resultsChecks.errorServices = false
	resultsChecks.errorServicesPlan = false

	fmt.Println("")
	fmt.Printf(InfoColor, "  Check plans")
	fmt.Println("")
	fmt.Printf(InfoColor, "-------------------------------------------------------------")
	fmt.Println("")

    for _, resource := range YAMLParsed.Resources {
		if resource.Type ==  c_TypeManaged{
			command_result, errorCliCommand := cliConnection.CliCommandWithoutTerminalOutput("marketplace", "-s", resource.Parameters.Service)
			if errorCliCommand == nil {
				errorServicesPlan = true
				// Parse marketplace services
				for i := 4; i < len(command_result); i++ {
					words := strings.Fields(command_result[i])
					if resource.Parameters.ServicePlan == words[0]{
						fmt.Printf(OkColor,resource.Name)
						if words[len(words)-1] == c_CostPaid{
							fmt.Printf(WarningColor," - [Paid]")
							fmt.Printf(" " + resource.Parameters.ServicePlan)
						}
						fmt.Println("")
						errorServicesPlan = false
						break
					}
				}
				if errorServicesPlan{
					fmt.Printf(ErrorColor,resource.Name)
					fmt.Printf(" - ")
					fmt.Printf(resource.Parameters.ServicePlan)
					fmt.Printf(" [service plan does not exist]");
					fmt.Println("")
					resultsChecks.errorServicesPlan = true
				}
			}else{
				fmt.Printf(ErrorColor,resource.Name )
				fmt.Printf(" - ")
				fmt.Printf(resource.Parameters.Service)
				fmt.Printf(" [service does not exist]")
				fmt.Println("")
				resultsChecks.errorServices = true
			}

		}
	}
}

// Check bindings
func (YAMLParsed YAMLData) CheckResourceListBinding(cliConnection plugin.CliConnection) {

	//var services plugin_models.GetService_Model
	var errorCliCommand error

	resultsChecks.errorBindings = false

	fmt.Println("")
	fmt.Printf(InfoColor, "  Check binding")
	fmt.Println("")
	fmt.Printf(InfoColor, "-------------------------------------------------------------")
	fmt.Println("")

	//Check servies
    for _, resource := range YAMLParsed.Resources {
		if resource.Type ==  c_TypeExisting{
			_,errorCliCommand = cliConnection.GetService(resource.Name)
			if errorCliCommand != nil {
				fmt.Printf(ErrorColor,resource.Name)
				fmt.Printf(" [service does not exist in your Org/Space] ")
				fmt.Println("")
				resultsChecks.errorBindings = true

			}else{
				fmt.Printf(OkColor,resource.Name)
				fmt.Println("")
			}
		}
	}
	
}

func (pluginDemo *PluginParams) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "check-before-deploy",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 1,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "check-before-deploy",
				Alias:    "cbd",
				HelpText: "Check the YAML file and services before deploying your MTA",
				UsageDetails: plugin.Usage{
					Usage: "cf check-before-deploy -file [path] or -mta [path] -check-binding -check-service -all",
					Options: map[string]string{
						"file": "Path with YAML file",
						"mta": "Path with MTA file - The 'META-INF/mtad.yaml' file must exist in the MTA",
						"check-binding": "Check the YAML file if the binded services exist in org / space",
						"check-service": "Check the YAML file if exist services plant exist in org / space",
						"all": "Active all validations",
					},
				},
			},
		},
	}
}