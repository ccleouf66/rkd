package main

import (
	"fmt"
	"rkd/helm"

	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

var (
	rancherStableRepoURL  = "https://releases.rancher.com/server-charts/stable"
	rancherLatestRepoURL  = "https://releases.rancher.com/server-charts/latest"
	rancherStableRepoName = "rancher-stable"
	rancherLatestRepoName = "rancher-latest"
	rancherChartName      = "rancher"

	testChartName = "ceph-csi-rbd"
	testRepoName  = "ceph-csi-rbd"
	testRepoURL   = "https://ceph.github.io/csi-charts"
)

type imageValue struct {
	name       string
	repository string
	tag        string
}

type k8sYamlStruct struct {
	APIVersion string          `json:"apiVersion" yaml:"apiVersion"`
	Kind       string          `json:"kind" yaml:"kind"`
	Metadata   k8sYamlMetadata `json:"metadata" yaml:"metadata"`
	Spec       k8sYamlSpec     `json:"spec" yaml:"spec"`
}

type k8sYamlMetadata struct {
	Namespace string
	Name      string
}

type k8sYamlSpec struct {
	Template k8sYamlTemplate
}

type k8sYamlTemplate struct {
	Spec k8sYamlSpecTemplate
}

type k8sYamlSpecTemplate struct {
	Containers []k8sYamlContrainer
}

type k8sYamlContrainer struct {
	Name  string
	Image string
}

func main() {

	///test
	helm.RepoAdd(testChartName, testRepoURL)
	helm.RepoUpdate()
	helm.DownloadChart(testRepoName, testChartName, "3.1.1", "./test")

	chart, err := loader.Load("./test/gitlab-4.4.4.tgz")
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	cvals, err := chartutil.CoalesceValues(chart, nil)
	if err != nil {
		fmt.Printf("ERR CoalesceValues => %s\n", err)
		return
	}

	getImage(cvals)

	// charts := chart.Dependencies()
	// charts = append(charts, chart)

	// fmt.Printf("%s\n", chart.Name())
	// fmt.Printf("%s\n", chart.AppVersion())

	// options := chartutil.ReleaseOptions{
	// 	Name:      "test-release",
	// 	Namespace: "default",
	// }

	// for _, c := range charts {

	// cValues := completMockData(c.Values)

	// fmt.Printf("%s", cValues)

	// valuesToRender, err := chartutil.ToRenderValues(c, cvals, options, nil)
	// if err != nil {
	// 	fmt.Printf("ERR ToRenderValues => %s\n", err)
	// 	return
	// }

	// var e engine.Engine
	// e.LintMode = true

	// renderedContentMap, err := e.Render(c, valuesToRender)
	// if err != nil {
	// 	fmt.Printf("ERR Render => %s\n", err)
	// 	return
	// }

	// for _, template := range c.Templates {
	// 	fileName := template.Name

	// 	// We only apply the following lint rules to yaml files
	// 	if filepath.Ext(fileName) != ".yaml" || filepath.Ext(fileName) == ".yml" {
	// 		continue
	// 	}

	// 	renderedContent := renderedContentMap[path.Join(c.Name(), fileName)]
	// 	if strings.TrimSpace(renderedContent) != "" {
	// 		var yamlStruct k8sYamlStruct
	// 		// Even though K8sYamlStruct only defines a few fields, an error in any other
	// 		// key will be raised as well
	// 		err := yaml.Unmarshal([]byte(renderedContent), &yamlStruct)
	// 		if err != nil {
	// 			fmt.Printf("ERR yaml Unmarshal => %s\n", err)
	// 			return
	// 		}

	// 		if !reflect.DeepEqual(yamlStruct.Spec, k8sYamlSpec{}) {
	// 			if !reflect.DeepEqual(yamlStruct.Spec.Template, k8sYamlTemplate{}) {
	// 				if !reflect.DeepEqual(yamlStruct.Spec.Template.Spec, k8sYamlSpecTemplate{}) {
	// 					if len(yamlStruct.Spec.Template.Spec.Containers) > 0 {
	// 						fmt.Printf("\n%s (%s)\n---\n", fileName, yamlStruct.Kind)
	// 						for _, c := range yamlStruct.Spec.Template.Spec.Containers {
	// 							fmt.Printf("%s,", c.Image)
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 		fmt.Printf("\n")
	// 	}
	// }
	// }

	/*
		for i, t := range chart.Templates {
			fmt.Printf("%d. %s\n", i, t.Name)

			if t.Name == "templates/deployment.yaml" {
				//fmt.Printf("%s\n", t.Data)
				d, err := chartutil.ReadValues(t.Data)
				if err != nil {
					fmt.Printf("%s\n", err)
				}
				fmt.Printf("%s\n", d)
			}
		}*/
	///

	// app := cli.NewApp()
	// app.Name = "rkd"
	// app.Usage = "Rancher Kubernetes Downloader"

	// app.Commands = []cli.Command{
	// 	cmd.ListCommand(),
	// 	cmd.DownloadCommand(),
	// }

	// app.Run(os.Args)
}

func getImage(values map[string]interface{}) {
	for _, v := range values {
		switch v.(type) {
		// case string:
		// 	if k == "image" || k == "tag" || k == "repository" {
		// 		fmt.Printf("%s => %s\n", k, v)
		// 	}
		case map[string]interface{}:
			var imgObj imageValue
			yamlObj, err := yaml.Marshal(v)
			if err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
			err = yaml.Unmarshal(yamlObj, &imgObj)
			if err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
			if imgObj.name != "" && imgObj.tag != "" {
				fmt.Printf("%s:%s\n", imgObj.name, imgObj.tag)
			} else if imgObj.repository != "" && imgObj.tag != "" {
				fmt.Printf("%s:%s\n", imgObj.repository, imgObj.tag)
			} else {
				fmt.Printf("%s\n", imgObj)
				getImage(v.(map[string]interface{}))
			}
		}
	}
}

func completMockData(obj map[string]interface{}) map[string]interface{} {
	for k, v := range obj {
		switch v.(type) {
		case bool:
			obj[k] = true
		case string:
			if obj[k] == "" || obj[k] == nil {
				obj[k] = "a"
			}
		case float64:
			if obj[k] == nil {
				obj[k] = 0
			}
		case int:
			if obj[k] == nil {
				obj[k] = 0
			}
		case map[string]interface{}:
			obj[k] = completMockData(v.(map[string]interface{}))
		case nil:
			obj[k] = true
		case []interface{}:
			obj[k] = completMockSliceData(v.([]interface{}))
		default:
			fmt.Printf("================================>    %T %s %s", v, k, v)
		}
	}
	return obj
}

func completMockSliceData(obj []interface{}) []interface{} {
	for k, v := range obj {
		switch v.(type) {
		case bool:
			obj[k] = false
		case string:
			if obj[k] == "" || obj[k] == nil {
				obj[k] = "a"
			}
		case float64:
			if obj[k] == nil {
				obj[k] = 0
			}
		case int:
			if obj[k] == nil {
				obj[k] = 0
			}
		case map[string]interface{}:
			obj[k] = completMockData(v.(map[string]interface{}))
		case nil:
			obj[k] = "a"
		default:
			fmt.Printf("================================>    %T", v)
		}
	}
	return obj
}
