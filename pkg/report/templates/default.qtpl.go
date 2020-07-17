// Code generated by qtc from "default.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line pkg/report/templates/default.qtpl:1
package templates

//line pkg/report/templates/default.qtpl:1
import "github.com/aquasecurity/starboard/pkg/apis/aquasecurity/v1alpha1"

//line pkg/report/templates/default.qtpl:2
import "github.com/aquasecurity/starboard/pkg/kube"

//line pkg/report/templates/default.qtpl:5
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line pkg/report/templates/default.qtpl:5
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line pkg/report/templates/default.qtpl:6
type ReportPage struct {
	VulnsReports       []v1alpha1.Vulnerability
	ConfigAuditReports []v1alpha1.ConfigAuditReport
	Workload           kube.Object
}

type ConfigAuditSummary struct {
	PodPass       int
	PodFail       int
	ContainerPass int
	ContainerFail int
}

//line pkg/report/templates/default.qtpl:20
func (p *ReportPage) StreamTitle(qw422016 *qt422016.Writer) {
//line pkg/report/templates/default.qtpl:20
	qw422016.N().S(`
	Starboard Security Report - `)
//line pkg/report/templates/default.qtpl:21
	qw422016.E().S(p.Workload.Namespace)
//line pkg/report/templates/default.qtpl:21
	qw422016.N().S(`/`)
//line pkg/report/templates/default.qtpl:21
	qw422016.E().S(string(p.Workload.Kind))
//line pkg/report/templates/default.qtpl:21
	qw422016.N().S(`/`)
//line pkg/report/templates/default.qtpl:21
	qw422016.E().S(p.Workload.Name)
//line pkg/report/templates/default.qtpl:21
	qw422016.N().S(`
`)
//line pkg/report/templates/default.qtpl:22
}

//line pkg/report/templates/default.qtpl:22
func (p *ReportPage) WriteTitle(qq422016 qtio422016.Writer) {
//line pkg/report/templates/default.qtpl:22
	qw422016 := qt422016.AcquireWriter(qq422016)
//line pkg/report/templates/default.qtpl:22
	p.StreamTitle(qw422016)
//line pkg/report/templates/default.qtpl:22
	qt422016.ReleaseWriter(qw422016)
//line pkg/report/templates/default.qtpl:22
}

//line pkg/report/templates/default.qtpl:22
func (p *ReportPage) Title() string {
//line pkg/report/templates/default.qtpl:22
	qb422016 := qt422016.AcquireByteBuffer()
//line pkg/report/templates/default.qtpl:22
	p.WriteTitle(qb422016)
//line pkg/report/templates/default.qtpl:22
	qs422016 := string(qb422016.B)
//line pkg/report/templates/default.qtpl:22
	qt422016.ReleaseByteBuffer(qb422016)
//line pkg/report/templates/default.qtpl:22
	return qs422016
//line pkg/report/templates/default.qtpl:22
}

//line pkg/report/templates/default.qtpl:26
func (p *ReportPage) GetMergedVulnsSummary() v1alpha1.VulnerabilitySummary {
	merged := v1alpha1.VulnerabilitySummary{}
	for _, report := range p.VulnsReports {
		merged.CriticalCount += report.Report.Summary.CriticalCount
		merged.HighCount += report.Report.Summary.HighCount
		merged.MediumCount += report.Report.Summary.MediumCount
		merged.LowCount += report.Report.Summary.LowCount
		merged.UnknownCount += report.Report.Summary.UnknownCount
	}
	return merged
}

//line pkg/report/templates/default.qtpl:40
func (p *ReportPage) GetConfigAuditSummary() ConfigAuditSummary {
	summary := ConfigAuditSummary{}
	// sum container checks
	for _, report := range p.ConfigAuditReports {
		for _, checks := range report.Report.ContainerChecks {
			for _, check := range checks {
				if check.Success {
					summary.ContainerPass += 1
				} else {
					summary.ContainerFail += 1
				}
			}
		}
		// sum pod checks
		for _, check := range report.Report.PodChecks {
			if check.Success {
				summary.PodPass += 1
			} else {
				summary.PodFail += 1
			}
		}
	}
	return summary
}

//line pkg/report/templates/default.qtpl:67
func (p *ReportPage) StreamBody(qw422016 *qt422016.Writer) {
//line pkg/report/templates/default.qtpl:67
	qw422016.N().S(`
	<style>
  a {
    color: inherit;
  }
  </style>
  <div class="container border-right border-left" style="height: 100%; overflow: scroll;">
            <div class="col mt-5">
                <div class="row text-center">
                    <img style="width: 25vh" class="mx-auto" src="https://www.aquasec.com/wp-content/uploads/2016/05/aqua_logo_fullcolor.png" alt="">
                </div>
                <div class="row mt-4 text-center">
                    <h2 class="text-muted mx-auto">Starboard Security Report for<br>
                </div>
                <div class="row text-center">
                    <h2 class="text-muted mx-auto">`)
//line pkg/report/templates/default.qtpl:82
	qw422016.E().S(string(p.Workload.Kind))
//line pkg/report/templates/default.qtpl:82
	qw422016.N().S(` `)
//line pkg/report/templates/default.qtpl:82
	qw422016.E().S(p.Workload.Name)
//line pkg/report/templates/default.qtpl:82
	qw422016.N().S(` in namespace `)
//line pkg/report/templates/default.qtpl:82
	qw422016.E().S(p.Workload.Namespace)
//line pkg/report/templates/default.qtpl:82
	qw422016.N().S(`</h2>
                </div>

		  
                <div class="row mt-5 px-3">
                    <h4>Table Of Contents</h4>
                </div>
                <div class="row">
                    <ul>
                        `)
//line pkg/report/templates/default.qtpl:91
	if len(p.VulnsReports) > 0 {
//line pkg/report/templates/default.qtpl:91
		qw422016.N().S(`
                        <li>
                            <a href="#vuln_header">Vulnerabilities</a></li>
                            <ul>
                              `)
//line pkg/report/templates/default.qtpl:95
		for count, report := range p.VulnsReports {
//line pkg/report/templates/default.qtpl:95
			qw422016.N().S(`
                                <li><a href="#vulns_container_`)
//line pkg/report/templates/default.qtpl:96
			qw422016.N().D(count)
//line pkg/report/templates/default.qtpl:96
			qw422016.N().S(`">`)
//line pkg/report/templates/default.qtpl:96
			qw422016.E().S(report.Labels["starboard.container.name"])
//line pkg/report/templates/default.qtpl:96
			qw422016.N().S(`</a></li>
                              `)
//line pkg/report/templates/default.qtpl:97
		}
//line pkg/report/templates/default.qtpl:97
		qw422016.N().S(`
                            </ul>
                        </li>
                        `)
//line pkg/report/templates/default.qtpl:100
	}
//line pkg/report/templates/default.qtpl:100
	qw422016.N().S(`
                        `)
//line pkg/report/templates/default.qtpl:101
	if len(p.ConfigAuditReports) > 0 {
//line pkg/report/templates/default.qtpl:101
		qw422016.N().S(`
                        <li>
                            <a href="#ca_header">Configuration Audit</a>
                            <ul>
                              <li><a href="#ca_pod_checks">Pod Checks</a></li>
                              `)
//line pkg/report/templates/default.qtpl:106
		for _, report := range p.ConfigAuditReports {
//line pkg/report/templates/default.qtpl:106
			qw422016.N().S(`
                                `)
//line pkg/report/templates/default.qtpl:107
			for container, _ := range report.Report.ContainerChecks {
//line pkg/report/templates/default.qtpl:107
				qw422016.N().S(`
                                  <li><a href="#ca_container_`)
//line pkg/report/templates/default.qtpl:108
				qw422016.E().S(container)
//line pkg/report/templates/default.qtpl:108
				qw422016.N().S(`">`)
//line pkg/report/templates/default.qtpl:108
				qw422016.E().S(container)
//line pkg/report/templates/default.qtpl:108
				qw422016.N().S(`</a></li>
                                `)
//line pkg/report/templates/default.qtpl:109
			}
//line pkg/report/templates/default.qtpl:109
			qw422016.N().S(`
                              `)
//line pkg/report/templates/default.qtpl:110
		}
//line pkg/report/templates/default.qtpl:110
		qw422016.N().S(`
                            </ul>
                        </li>
                        `)
//line pkg/report/templates/default.qtpl:113
	}
//line pkg/report/templates/default.qtpl:113
	qw422016.N().S(`
                    </ul>
                </div>


                `)
//line pkg/report/templates/default.qtpl:118
	if len(p.VulnsReports) > 0 {
//line pkg/report/templates/default.qtpl:118
		qw422016.N().S(`
                <!-- Vulnerabilities -->
                <div class="row text-center border-bottom mt-4">
                    <h3 class="mx-auto " id="vuln_header" style="color: rgb(0, 160, 170);">Vulnerabilities</h3>
                </div>
                <!-- Cards -->
                <div class="">
                    <div class="row my-5" style="font-size:small;">
                        <!-- Scanner -->
                        <div class="col-3 border rounded shadow px-3 py-2 ml-4 ">
                            <div class="row text-center">
                                <div class="col">
                                    <p class="mb-2 pb-1 border-bottom">Scanner</p>
                                </div>
                             </div>
                             <div class="row">
                                <div class="col">
                                    <p class="my-0">Name:  `)
//line pkg/report/templates/default.qtpl:135
		qw422016.E().S(p.VulnsReports[0].Report.Scanner.Name)
//line pkg/report/templates/default.qtpl:135
		qw422016.N().S(`</p>
                                    <p class="my-0">Vendor:  `)
//line pkg/report/templates/default.qtpl:136
		qw422016.E().S(p.VulnsReports[0].Report.Scanner.Vendor)
//line pkg/report/templates/default.qtpl:136
		qw422016.N().S(`</p>
                                    <p class="my-0">Version:  `)
//line pkg/report/templates/default.qtpl:137
		qw422016.E().S(p.VulnsReports[0].Report.Scanner.Version)
//line pkg/report/templates/default.qtpl:137
		qw422016.N().S(`</p>
                                </div>
                             </div>
                        </div>
                        <!-- summary -->
                        <div class="col-5 border rounded shadow py-2 mx-auto ">
                            <div class="row text-center">
                               <div class="col">
                                   <p class="mb-2 pb-1 border-bottom">Summary</p>
                               </div>
                            </div>
                            <div class="row">
                                `)
//line pkg/report/templates/default.qtpl:150
		summary := p.GetMergedVulnsSummary()

//line pkg/report/templates/default.qtpl:151
		qw422016.N().S(`
                                `)
//line pkg/report/templates/default.qtpl:152
		if summary.CriticalCount > 0 {
//line pkg/report/templates/default.qtpl:152
			qw422016.N().S(`
                                <div class="col text-center p-0 text-danger font-weight-bold">
                                `)
//line pkg/report/templates/default.qtpl:154
		} else {
//line pkg/report/templates/default.qtpl:154
			qw422016.N().S(`
                                <div class="col text-center p-0">
                                `)
//line pkg/report/templates/default.qtpl:156
		}
//line pkg/report/templates/default.qtpl:156
		qw422016.N().S(`
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:157
		qw422016.N().D(summary.CriticalCount)
//line pkg/report/templates/default.qtpl:157
		qw422016.N().S(`</p>
                                    <p class="mx-auto ">CRITICAL</p>
                                </div>
                                `)
//line pkg/report/templates/default.qtpl:160
		if summary.HighCount > 0 {
//line pkg/report/templates/default.qtpl:160
			qw422016.N().S(`
                                <div class="col text-center p-0 text-danger font-weight-bold">
                                `)
//line pkg/report/templates/default.qtpl:162
		} else {
//line pkg/report/templates/default.qtpl:162
			qw422016.N().S(`
                                <div class="col text-center p-0">
                                `)
//line pkg/report/templates/default.qtpl:164
		}
//line pkg/report/templates/default.qtpl:164
		qw422016.N().S(`
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:165
		qw422016.N().D(summary.HighCount)
//line pkg/report/templates/default.qtpl:165
		qw422016.N().S(`</p>
                                    <p class="mx-auto ">HIGH</p>
                                </div>
                                `)
//line pkg/report/templates/default.qtpl:168
		if summary.MediumCount > 0 {
//line pkg/report/templates/default.qtpl:168
			qw422016.N().S(`
                                <div class="col text-center p-0 text-warning font-weight-bold">
                                `)
//line pkg/report/templates/default.qtpl:170
		} else {
//line pkg/report/templates/default.qtpl:170
			qw422016.N().S(`
                                <div class="col text-center p-0">
                                `)
//line pkg/report/templates/default.qtpl:172
		}
//line pkg/report/templates/default.qtpl:172
		qw422016.N().S(`
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:173
		qw422016.N().D(summary.MediumCount)
//line pkg/report/templates/default.qtpl:173
		qw422016.N().S(`</p>
                                    <p class="mx-auto ">MEDIUM</p>
                                </div>
                                <div class="col text-center p-0">
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:177
		qw422016.N().D(summary.LowCount)
//line pkg/report/templates/default.qtpl:177
		qw422016.N().S(`</p>
                                    <p class="mx-auto ">LOW</p>
                                </div>
                                <div class="col text-center p-0">
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:181
		qw422016.N().D(summary.UnknownCount)
//line pkg/report/templates/default.qtpl:181
		qw422016.N().S(`</p>
                                    <p class="mx-auto ">UNKNOWN</p>
                                </div>
                            </div>
                        </div>
                        <!-- Metadata -->
                        <div class="col-3 border rounded shadow px-3 py-2 mr-4">
                            <div class="row text-center">
                                <div class="col">
                                    <p class="mb-2 pb-1 border-bottom">Metadata</p>
                                </div>
                             </div>
                             <div class="row">
                                <div class="col">
                                    <p class="my-0">Generated at:  `)
//line pkg/report/templates/default.qtpl:195
		qw422016.E().S(p.VulnsReports[0].CreationTimestamp.String())
//line pkg/report/templates/default.qtpl:195
		qw422016.N().S(`</p>
                                </div>
                             </div>
                        </div>

                    </div>      
                </div>
                `)
//line pkg/report/templates/default.qtpl:202
	}
//line pkg/report/templates/default.qtpl:202
	qw422016.N().S(`
                
                `)
//line pkg/report/templates/default.qtpl:204
	for count, report := range p.VulnsReports {
//line pkg/report/templates/default.qtpl:204
		qw422016.N().S(`
                
                  <div class="row"><h5 class="text-info" id="vulns_container_`)
//line pkg/report/templates/default.qtpl:206
		qw422016.N().D(count)
//line pkg/report/templates/default.qtpl:206
		qw422016.N().S(`">Container `)
//line pkg/report/templates/default.qtpl:206
		qw422016.E().S(report.Labels["starboard.container.name"])
//line pkg/report/templates/default.qtpl:206
		qw422016.N().S(`</h5></div>
                  <div class="row"><p>`)
//line pkg/report/templates/default.qtpl:207
		qw422016.E().S(report.Report.Registry.URL)
//line pkg/report/templates/default.qtpl:207
		qw422016.N().S(`/`)
//line pkg/report/templates/default.qtpl:207
		qw422016.E().S(report.Report.Artifact.Repository)
//line pkg/report/templates/default.qtpl:207
		qw422016.N().S(`:`)
//line pkg/report/templates/default.qtpl:207
		qw422016.E().S(report.Report.Artifact.Tag)
//line pkg/report/templates/default.qtpl:207
		qw422016.N().S(`</p></div>
                  `)
//line pkg/report/templates/default.qtpl:208
		if len(report.Report.Vulnerabilities) == 0 {
//line pkg/report/templates/default.qtpl:208
			qw422016.N().S(`
                    <div class="row">
                      <p class="alert alert-success py-0 m-0" style="font-size: small;">No Vulnerabilities</p>
                    </div>                  
                  `)
//line pkg/report/templates/default.qtpl:212
		} else {
//line pkg/report/templates/default.qtpl:212
			qw422016.N().S(`

                  <div class="row">
                      <table class="table table-sm table-bordered">
                          <thead>
                              <tr>
                                <th scope="col">ID</th>
                                <th scope="col">Severity</th>
                                <th scope="col">Resource</th>
                                <th scope="col">Installed Version</th>
                                <th scope="col">Fixed Version</th>
                              </tr>
                            </thead>
                            <tbody>
                  `)
//line pkg/report/templates/default.qtpl:226
			for _, v := range report.Report.Vulnerabilities {
//line pkg/report/templates/default.qtpl:226
				qw422016.N().S(`
                    <tr>
                      <td>`)
//line pkg/report/templates/default.qtpl:228
				qw422016.E().S(v.VulnerabilityID)
//line pkg/report/templates/default.qtpl:228
				qw422016.N().S(`</td>
                      <td>`)
//line pkg/report/templates/default.qtpl:229
				qw422016.E().S(string(v.Severity))
//line pkg/report/templates/default.qtpl:229
				qw422016.N().S(`</td>
                      <td>`)
//line pkg/report/templates/default.qtpl:230
				qw422016.E().S(v.Resource)
//line pkg/report/templates/default.qtpl:230
				qw422016.N().S(`</td>
                      <td>`)
//line pkg/report/templates/default.qtpl:231
				qw422016.E().S(v.InstalledVersion)
//line pkg/report/templates/default.qtpl:231
				qw422016.N().S(`</td>
                      <td>`)
//line pkg/report/templates/default.qtpl:232
				qw422016.E().S(v.FixedVersion)
//line pkg/report/templates/default.qtpl:232
				qw422016.N().S(`</td>
                    </tr>	
                  `)
//line pkg/report/templates/default.qtpl:234
			}
//line pkg/report/templates/default.qtpl:234
			qw422016.N().S(`
                            </tbody>
                      </table>
                  </div>
                `)
//line pkg/report/templates/default.qtpl:238
		}
//line pkg/report/templates/default.qtpl:238
		qw422016.N().S(`
                `)
//line pkg/report/templates/default.qtpl:239
	}
//line pkg/report/templates/default.qtpl:239
	qw422016.N().S(`
                

                <!-- Config Audits -->
                `)
//line pkg/report/templates/default.qtpl:243
	if len(p.ConfigAuditReports) > 0 {
//line pkg/report/templates/default.qtpl:243
		qw422016.N().S(`
                  <div class="row pt-3 text-center border-bottom my-4">
                      <h3 class="mx-auto" id="ca_header" style="color: rgb(0, 160, 170);">Configuration Audit</h3>
                  </div>
                  <!-- Cards -->
                <div class="">
                    <div class="row my-5" style="font-size:small;">
                        <!-- Scanner -->
                        <div class="col-3 border rounded shadow px-3 py-2 ml-4 ">
                            <div class="row text-center">
                                <div class="col">
                                    <p class="mb-2 pb-1 border-bottom">Scanner</p>
                                </div>
                             </div>
                             <div class="row">
                                <div class="col">
                                    <p class="my-0">Name:  `)
//line pkg/report/templates/default.qtpl:259
		qw422016.E().S(p.ConfigAuditReports[0].Report.Scanner.Name)
//line pkg/report/templates/default.qtpl:259
		qw422016.N().S(`</p>
                                    <p class="my-0">Vendor:  `)
//line pkg/report/templates/default.qtpl:260
		qw422016.E().S(p.ConfigAuditReports[0].Report.Scanner.Vendor)
//line pkg/report/templates/default.qtpl:260
		qw422016.N().S(`</p>
                                    <p class="my-0">Version:  `)
//line pkg/report/templates/default.qtpl:261
		qw422016.E().S(p.ConfigAuditReports[0].Report.Scanner.Version)
//line pkg/report/templates/default.qtpl:261
		qw422016.N().S(`</p>
                                </div>
                             </div>
                        </div>
                        <!-- summary -->
                        <div class="col-3 border rounded shadow py-2 mx-auto ">
                            <div class="row text-center">
                               <div class="col">
                                   <p class="mb-2 pb-1 border-bottom">Summary</p>
                               </div>
                            </div>
                            <div class="row">
                                `)
//line pkg/report/templates/default.qtpl:274
		summary := p.GetConfigAuditSummary()
		sumPass := summary.PodPass + summary.ContainerPass
		sumFail := summary.PodFail + summary.ContainerFail

//line pkg/report/templates/default.qtpl:277
		qw422016.N().S(`
                                `)
//line pkg/report/templates/default.qtpl:278
		if sumPass > 0 {
//line pkg/report/templates/default.qtpl:278
			qw422016.N().S(`
                                <div class="col text-center p-0 text-success font-weight-bold">
                                `)
//line pkg/report/templates/default.qtpl:280
		} else {
//line pkg/report/templates/default.qtpl:280
			qw422016.N().S(`
                                <div class="col text-center p-0">
                                `)
//line pkg/report/templates/default.qtpl:282
		}
//line pkg/report/templates/default.qtpl:282
		qw422016.N().S(`
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:283
		qw422016.N().D(sumPass)
//line pkg/report/templates/default.qtpl:283
		qw422016.N().S(`</p>
                                    <p class="mx-auto ">PASS</p>
                                </div>
                                `)
//line pkg/report/templates/default.qtpl:286
		if sumFail > 0 {
//line pkg/report/templates/default.qtpl:286
			qw422016.N().S(`
                                <div class="col text-center p-0 text-danger font-weight-bold">
                                `)
//line pkg/report/templates/default.qtpl:288
		} else {
//line pkg/report/templates/default.qtpl:288
			qw422016.N().S(`
                                <div class="col text-center p-0">
                                `)
//line pkg/report/templates/default.qtpl:290
		}
//line pkg/report/templates/default.qtpl:290
		qw422016.N().S(`
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:291
		qw422016.N().D(sumFail)
//line pkg/report/templates/default.qtpl:291
		qw422016.N().S(`</p>
                                    <p class="mx-auto ">FAIL</p>
                                </div>
                            </div>
                        </div>
                        <!-- Metadata -->
                        <div class="col-3 border rounded shadow px-3 py-2 mr-4">
                            <div class="row text-center">
                                <div class="col">
                                    <p class="mb-2 pb-1 border-bottom">Metadata</p>
                                </div>
                             </div>
                             <div class="row">
                                <div class="col">
                                    <p class="my-0">Generated at:  `)
//line pkg/report/templates/default.qtpl:305
		qw422016.E().S(p.ConfigAuditReports[0].CreationTimestamp.String())
//line pkg/report/templates/default.qtpl:305
		qw422016.N().S(`</p>
                                </div>
                             </div>
                        </div>

                    </div>      
                </div>
                `)
//line pkg/report/templates/default.qtpl:312
	}
//line pkg/report/templates/default.qtpl:312
	qw422016.N().S(`
                `)
//line pkg/report/templates/default.qtpl:313
	for _, report := range p.ConfigAuditReports {
//line pkg/report/templates/default.qtpl:313
		qw422016.N().S(`        
                  <div class="row"><h5 class="text-info" id="ca_pod_checks">Pod Checks</h5></div>
                  <div class="row">
                      <table class="table table-sm table-bordered">
                          <thead>
                              <tr>
                                <th scope="col">Success</th>
                                <th scope="col">ID</th>
                                <th scope="col">Severity</th>
                                <th scope="col">Category</th>
                              </tr>
                            </thead>
                            <tbody>
                              `)
//line pkg/report/templates/default.qtpl:326
		for _, check := range report.Report.PodChecks {
//line pkg/report/templates/default.qtpl:326
			qw422016.N().S(`
                                <tr>
                                  <td>`)
//line pkg/report/templates/default.qtpl:328
			qw422016.E().V(check.Success)
//line pkg/report/templates/default.qtpl:328
			qw422016.N().S(`</td>
                                  <td>`)
//line pkg/report/templates/default.qtpl:329
			qw422016.E().S(check.ID)
//line pkg/report/templates/default.qtpl:329
			qw422016.N().S(`</td>
                                  <td>`)
//line pkg/report/templates/default.qtpl:330
			qw422016.E().S(check.Severity)
//line pkg/report/templates/default.qtpl:330
			qw422016.N().S(`</td>
                                  <td>`)
//line pkg/report/templates/default.qtpl:331
			qw422016.E().S(check.Category)
//line pkg/report/templates/default.qtpl:331
			qw422016.N().S(`</td>
                                </tr>
                              `)
//line pkg/report/templates/default.qtpl:333
		}
//line pkg/report/templates/default.qtpl:333
		qw422016.N().S(`
                            </tbody>
                      </table>
                  </div>
                  `)
//line pkg/report/templates/default.qtpl:337
		for container, checks := range report.Report.ContainerChecks {
//line pkg/report/templates/default.qtpl:337
			qw422016.N().S(`
                    <div class="row"><h5 class="text-info" id="ca_container_`)
//line pkg/report/templates/default.qtpl:338
			qw422016.E().S(container)
//line pkg/report/templates/default.qtpl:338
			qw422016.N().S(`">Container `)
//line pkg/report/templates/default.qtpl:338
			qw422016.E().S(container)
//line pkg/report/templates/default.qtpl:338
			qw422016.N().S(`</h5></div>
                    <div class="row">
                        <table class="table table-sm table-bordered">
                            <thead>
                                <tr>
                                  <th scope="col">Success</th>
                                  <th scope="col">ID</th>
                                  <th scope="col">Severity</th>
                                  <th scope="col">Category</th>
                                </tr>
                              </thead>
                              <tbody>
                                `)
//line pkg/report/templates/default.qtpl:350
			for _, check := range checks {
//line pkg/report/templates/default.qtpl:350
				qw422016.N().S(`
                                  <tr>
                                    <td>`)
//line pkg/report/templates/default.qtpl:352
				qw422016.E().V(check.Success)
//line pkg/report/templates/default.qtpl:352
				qw422016.N().S(`</td>
                                    <td>`)
//line pkg/report/templates/default.qtpl:353
				qw422016.E().S(check.ID)
//line pkg/report/templates/default.qtpl:353
				qw422016.N().S(`</td>
                                    <td>`)
//line pkg/report/templates/default.qtpl:354
				qw422016.E().S(check.Severity)
//line pkg/report/templates/default.qtpl:354
				qw422016.N().S(`</td>
                                    <td>`)
//line pkg/report/templates/default.qtpl:355
				qw422016.E().S(check.Category)
//line pkg/report/templates/default.qtpl:355
				qw422016.N().S(`</td>
                                  </tr>
                                `)
//line pkg/report/templates/default.qtpl:357
			}
//line pkg/report/templates/default.qtpl:357
			qw422016.N().S(`
                              </tbody>
                        </table>
                    </div>
                  `)
//line pkg/report/templates/default.qtpl:361
		}
//line pkg/report/templates/default.qtpl:361
		qw422016.N().S(`
                `)
//line pkg/report/templates/default.qtpl:362
	}
//line pkg/report/templates/default.qtpl:362
	qw422016.N().S(`
            </div>
        </div>
`)
//line pkg/report/templates/default.qtpl:365
}

//line pkg/report/templates/default.qtpl:365
func (p *ReportPage) WriteBody(qq422016 qtio422016.Writer) {
//line pkg/report/templates/default.qtpl:365
	qw422016 := qt422016.AcquireWriter(qq422016)
//line pkg/report/templates/default.qtpl:365
	p.StreamBody(qw422016)
//line pkg/report/templates/default.qtpl:365
	qt422016.ReleaseWriter(qw422016)
//line pkg/report/templates/default.qtpl:365
}

//line pkg/report/templates/default.qtpl:365
func (p *ReportPage) Body() string {
//line pkg/report/templates/default.qtpl:365
	qb422016 := qt422016.AcquireByteBuffer()
//line pkg/report/templates/default.qtpl:365
	p.WriteBody(qb422016)
//line pkg/report/templates/default.qtpl:365
	qs422016 := string(qb422016.B)
//line pkg/report/templates/default.qtpl:365
	qt422016.ReleaseByteBuffer(qb422016)
//line pkg/report/templates/default.qtpl:365
	return qs422016
//line pkg/report/templates/default.qtpl:365
}
