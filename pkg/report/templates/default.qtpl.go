// Code generated by qtc from "default.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line pkg/report/templates/default.qtpl:1
package templates

//line pkg/report/templates/default.qtpl:1
import "github.com/aquasecurity/starboard/pkg/apis/aquasecurity/v1alpha1"

//line pkg/report/templates/default.qtpl:2
import "github.com/aquasecurity/starboard/pkg/kube"

//line pkg/report/templates/default.qtpl:3
import "github.com/aquasecurity/starboard/pkg/vulnerabilityreport"

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
	VulnsReports      vulnerabilityreport.WorkloadVulnerabilities
	ConfigAuditReport *v1alpha1.ConfigAuditReport
	Workload          kube.Object
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
		merged.CriticalCount += report.Summary.CriticalCount
		merged.HighCount += report.Summary.HighCount
		merged.MediumCount += report.Summary.MediumCount
		merged.LowCount += report.Summary.LowCount
		merged.UnknownCount += report.Summary.UnknownCount
	}
	return merged
}

//line pkg/report/templates/default.qtpl:40
func (p *ReportPage) GetConfigAuditSummary() ConfigAuditSummary {
	summary := ConfigAuditSummary{}
	// sum container checks
	for _, checks := range p.ConfigAuditReport.Report.ContainerChecks {
		for _, check := range checks {
			if check.Success {
				summary.ContainerPass += 1
			} else {
				summary.ContainerFail += 1
			}
		}
	}
	// sum pod checks
	for _, check := range p.ConfigAuditReport.Report.PodChecks {
		if check.Success {
			summary.PodPass += 1
		} else {
			summary.PodFail += 1
		}
	}

	return summary
}

//line pkg/report/templates/default.qtpl:66
func (p *ReportPage) StreamBody(qw422016 *qt422016.Writer) {
//line pkg/report/templates/default.qtpl:66
	qw422016.N().S(`
	<style>
  a {
    color: inherit;
  }
  @media print {
    .container {
        display: inline;
    }
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
//line pkg/report/templates/default.qtpl:86
	qw422016.E().S(string(p.Workload.Kind))
//line pkg/report/templates/default.qtpl:86
	qw422016.N().S(` `)
//line pkg/report/templates/default.qtpl:86
	qw422016.E().S(p.Workload.Name)
//line pkg/report/templates/default.qtpl:86
	qw422016.N().S(` in namespace `)
//line pkg/report/templates/default.qtpl:86
	qw422016.E().S(p.Workload.Namespace)
//line pkg/report/templates/default.qtpl:86
	qw422016.N().S(`</h2>
                </div>

		  
                <div class="row mt-5 px-3">
                    <h4>Table Of Contents</h4>
                </div>
                <div class="row">
                    <ul>
                        `)
//line pkg/report/templates/default.qtpl:95
	if len(p.VulnsReports) > 0 {
//line pkg/report/templates/default.qtpl:95
		qw422016.N().S(`
                        <li>
                            <a href="#vuln_header">Vulnerabilities</a></li>
                            <ul>
                              `)
//line pkg/report/templates/default.qtpl:99
		for container, _ := range p.VulnsReports {
//line pkg/report/templates/default.qtpl:99
			qw422016.N().S(`
                                <li><a href="#vulns_container_`)
//line pkg/report/templates/default.qtpl:100
			qw422016.E().S(container)
//line pkg/report/templates/default.qtpl:100
			qw422016.N().S(`">`)
//line pkg/report/templates/default.qtpl:100
			qw422016.E().S(container)
//line pkg/report/templates/default.qtpl:100
			qw422016.N().S(`</a></li>
                              `)
//line pkg/report/templates/default.qtpl:101
		}
//line pkg/report/templates/default.qtpl:101
		qw422016.N().S(`
                            </ul>
                        </li>
                        `)
//line pkg/report/templates/default.qtpl:104
	}
//line pkg/report/templates/default.qtpl:104
	qw422016.N().S(`
                        `)
//line pkg/report/templates/default.qtpl:105
	if p.ConfigAuditReport != nil && len(p.ConfigAuditReport.Report.PodChecks) > 0 {
//line pkg/report/templates/default.qtpl:105
		qw422016.N().S(`
                        <li>
                            <a href="#ca_header">Configuration Audit</a>
                            <ul>
                              <li><a href="#ca_pod_checks">Pod Checks</a></li>
                                `)
//line pkg/report/templates/default.qtpl:110
		for container, _ := range p.ConfigAuditReport.Report.ContainerChecks {
//line pkg/report/templates/default.qtpl:110
			qw422016.N().S(`
                                  <li><a href="#ca_container_`)
//line pkg/report/templates/default.qtpl:111
			qw422016.E().S(container)
//line pkg/report/templates/default.qtpl:111
			qw422016.N().S(`">`)
//line pkg/report/templates/default.qtpl:111
			qw422016.E().S(container)
//line pkg/report/templates/default.qtpl:111
			qw422016.N().S(`</a></li>
                                `)
//line pkg/report/templates/default.qtpl:112
		}
//line pkg/report/templates/default.qtpl:112
		qw422016.N().S(`
                            </ul>
                        </li>
                        `)
//line pkg/report/templates/default.qtpl:115
	}
//line pkg/report/templates/default.qtpl:115
	qw422016.N().S(`
                    </ul>
                </div>


                `)
//line pkg/report/templates/default.qtpl:120
	if len(p.VulnsReports) > 0 {
//line pkg/report/templates/default.qtpl:120
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
                                `)
//line pkg/report/templates/default.qtpl:138
		var scanner_name, scanner_vendor, scanner_version, creation_timestamp string
		for _, report := range p.VulnsReports {
			scanner_name = report.Scanner.Name
			scanner_vendor = report.Scanner.Vendor
			scanner_version = report.Scanner.Version
			creation_timestamp = report.UpdateTimestamp.String()
			break
		}

//line pkg/report/templates/default.qtpl:146
		qw422016.N().S(`
                                    <p class="my-0">Name:  `)
//line pkg/report/templates/default.qtpl:147
		qw422016.E().S(scanner_name)
//line pkg/report/templates/default.qtpl:147
		qw422016.N().S(`</p>
                                    <p class="my-0">Vendor:  `)
//line pkg/report/templates/default.qtpl:148
		qw422016.E().S(scanner_vendor)
//line pkg/report/templates/default.qtpl:148
		qw422016.N().S(`</p>
                                    <p class="my-0">Version:  `)
//line pkg/report/templates/default.qtpl:149
		qw422016.E().S(scanner_version)
//line pkg/report/templates/default.qtpl:149
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
//line pkg/report/templates/default.qtpl:162
		summary := p.GetMergedVulnsSummary()

//line pkg/report/templates/default.qtpl:163
		qw422016.N().S(`
                                `)
//line pkg/report/templates/default.qtpl:164
		if summary.CriticalCount > 0 {
//line pkg/report/templates/default.qtpl:164
			qw422016.N().S(`
                                <div class="col text-center p-0 text-danger font-weight-bold">
                                `)
//line pkg/report/templates/default.qtpl:166
		} else {
//line pkg/report/templates/default.qtpl:166
			qw422016.N().S(`
                                <div class="col text-center p-0">
                                `)
//line pkg/report/templates/default.qtpl:168
		}
//line pkg/report/templates/default.qtpl:168
		qw422016.N().S(`
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:169
		qw422016.N().D(summary.CriticalCount)
//line pkg/report/templates/default.qtpl:169
		qw422016.N().S(`</p>
                                    <p class="mx-auto ">CRITICAL</p>
                                </div>
                                `)
//line pkg/report/templates/default.qtpl:172
		if summary.HighCount > 0 {
//line pkg/report/templates/default.qtpl:172
			qw422016.N().S(`
                                <div class="col text-center p-0 text-danger font-weight-bold">
                                `)
//line pkg/report/templates/default.qtpl:174
		} else {
//line pkg/report/templates/default.qtpl:174
			qw422016.N().S(`
                                <div class="col text-center p-0">
                                `)
//line pkg/report/templates/default.qtpl:176
		}
//line pkg/report/templates/default.qtpl:176
		qw422016.N().S(`
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:177
		qw422016.N().D(summary.HighCount)
//line pkg/report/templates/default.qtpl:177
		qw422016.N().S(`</p>
                                    <p class="mx-auto ">HIGH</p>
                                </div>
                                `)
//line pkg/report/templates/default.qtpl:180
		if summary.MediumCount > 0 {
//line pkg/report/templates/default.qtpl:180
			qw422016.N().S(`
                                <div class="col text-center p-0 text-warning font-weight-bold">
                                `)
//line pkg/report/templates/default.qtpl:182
		} else {
//line pkg/report/templates/default.qtpl:182
			qw422016.N().S(`
                                <div class="col text-center p-0">
                                `)
//line pkg/report/templates/default.qtpl:184
		}
//line pkg/report/templates/default.qtpl:184
		qw422016.N().S(`
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:185
		qw422016.N().D(summary.MediumCount)
//line pkg/report/templates/default.qtpl:185
		qw422016.N().S(`</p>
                                    <p class="mx-auto ">MEDIUM</p>
                                </div>
                                <div class="col text-center p-0">
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:189
		qw422016.N().D(summary.LowCount)
//line pkg/report/templates/default.qtpl:189
		qw422016.N().S(`</p>
                                    <p class="mx-auto ">LOW</p>
                                </div>
                                <div class="col text-center p-0">
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:193
		qw422016.N().D(summary.UnknownCount)
//line pkg/report/templates/default.qtpl:193
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
//line pkg/report/templates/default.qtpl:207
		qw422016.E().S(creation_timestamp)
//line pkg/report/templates/default.qtpl:207
		qw422016.N().S(`</p>
                                </div>
                             </div>
                        </div>

                    </div>      
                </div>
                `)
//line pkg/report/templates/default.qtpl:214
	}
//line pkg/report/templates/default.qtpl:214
	qw422016.N().S(`
                
                `)
//line pkg/report/templates/default.qtpl:216
	for container, report := range p.VulnsReports {
//line pkg/report/templates/default.qtpl:216
		qw422016.N().S(`
                
                  <div class="row"><h5 class="text-info" id="vulns_container_`)
//line pkg/report/templates/default.qtpl:218
		qw422016.E().S(container)
//line pkg/report/templates/default.qtpl:218
		qw422016.N().S(`">Container `)
//line pkg/report/templates/default.qtpl:218
		qw422016.E().S(container)
//line pkg/report/templates/default.qtpl:218
		qw422016.N().S(`</h5></div>
                  <div class="row"><p>`)
//line pkg/report/templates/default.qtpl:219
		qw422016.E().S(report.Registry.Server)
//line pkg/report/templates/default.qtpl:219
		qw422016.N().S(`/`)
//line pkg/report/templates/default.qtpl:219
		qw422016.E().S(report.Artifact.Repository)
//line pkg/report/templates/default.qtpl:219
		qw422016.N().S(`:`)
//line pkg/report/templates/default.qtpl:219
		qw422016.E().S(report.Artifact.Tag)
//line pkg/report/templates/default.qtpl:219
		qw422016.N().S(`</p></div>
                  `)
//line pkg/report/templates/default.qtpl:220
		if len(report.Vulnerabilities) == 0 {
//line pkg/report/templates/default.qtpl:220
			qw422016.N().S(`
                    <div class="row">
                      <p class="alert alert-success py-0 m-0" style="font-size: small;">No Vulnerabilities</p>
                    </div>                  
                  `)
//line pkg/report/templates/default.qtpl:224
		} else {
//line pkg/report/templates/default.qtpl:224
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
//line pkg/report/templates/default.qtpl:238
			for _, v := range report.Vulnerabilities {
//line pkg/report/templates/default.qtpl:238
				qw422016.N().S(`
                    <tr>
                      <td>`)
//line pkg/report/templates/default.qtpl:240
				qw422016.E().S(v.VulnerabilityID)
//line pkg/report/templates/default.qtpl:240
				qw422016.N().S(`</td>
                      <td>`)
//line pkg/report/templates/default.qtpl:241
				qw422016.E().S(string(v.Severity))
//line pkg/report/templates/default.qtpl:241
				qw422016.N().S(`</td>
                      <td>`)
//line pkg/report/templates/default.qtpl:242
				qw422016.E().S(v.Resource)
//line pkg/report/templates/default.qtpl:242
				qw422016.N().S(`</td>
                      <td>`)
//line pkg/report/templates/default.qtpl:243
				qw422016.E().S(v.InstalledVersion)
//line pkg/report/templates/default.qtpl:243
				qw422016.N().S(`</td>
                      <td>`)
//line pkg/report/templates/default.qtpl:244
				qw422016.E().S(v.FixedVersion)
//line pkg/report/templates/default.qtpl:244
				qw422016.N().S(`</td>
                    </tr>	
                  `)
//line pkg/report/templates/default.qtpl:246
			}
//line pkg/report/templates/default.qtpl:246
			qw422016.N().S(`
                            </tbody>
                      </table>
                  </div>
                `)
//line pkg/report/templates/default.qtpl:250
		}
//line pkg/report/templates/default.qtpl:250
		qw422016.N().S(`
                `)
//line pkg/report/templates/default.qtpl:251
	}
//line pkg/report/templates/default.qtpl:251
	qw422016.N().S(`
                

                <!-- Config Audits -->
                `)
//line pkg/report/templates/default.qtpl:255
	if p.ConfigAuditReport != nil && len(p.ConfigAuditReport.Report.PodChecks) > 0 {
//line pkg/report/templates/default.qtpl:255
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
//line pkg/report/templates/default.qtpl:271
		qw422016.E().S(p.ConfigAuditReport.Report.Scanner.Name)
//line pkg/report/templates/default.qtpl:271
		qw422016.N().S(`</p>
                                    <p class="my-0">Vendor:  `)
//line pkg/report/templates/default.qtpl:272
		qw422016.E().S(p.ConfigAuditReport.Report.Scanner.Vendor)
//line pkg/report/templates/default.qtpl:272
		qw422016.N().S(`</p>
                                    <p class="my-0">Version:  `)
//line pkg/report/templates/default.qtpl:273
		qw422016.E().S(p.ConfigAuditReport.Report.Scanner.Version)
//line pkg/report/templates/default.qtpl:273
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
//line pkg/report/templates/default.qtpl:286
		summary := p.GetConfigAuditSummary()
		sumPass := summary.PodPass + summary.ContainerPass
		sumFail := summary.PodFail + summary.ContainerFail

//line pkg/report/templates/default.qtpl:289
		qw422016.N().S(`
                                `)
//line pkg/report/templates/default.qtpl:290
		if sumPass > 0 {
//line pkg/report/templates/default.qtpl:290
			qw422016.N().S(`
                                <div class="col text-center p-0 text-success font-weight-bold">
                                `)
//line pkg/report/templates/default.qtpl:292
		} else {
//line pkg/report/templates/default.qtpl:292
			qw422016.N().S(`
                                <div class="col text-center p-0">
                                `)
//line pkg/report/templates/default.qtpl:294
		}
//line pkg/report/templates/default.qtpl:294
		qw422016.N().S(`
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:295
		qw422016.N().D(sumPass)
//line pkg/report/templates/default.qtpl:295
		qw422016.N().S(`</p>
                                    <p class="mx-auto ">PASS</p>
                                </div>
                                `)
//line pkg/report/templates/default.qtpl:298
		if sumFail > 0 {
//line pkg/report/templates/default.qtpl:298
			qw422016.N().S(`
                                <div class="col text-center p-0 text-danger font-weight-bold">
                                `)
//line pkg/report/templates/default.qtpl:300
		} else {
//line pkg/report/templates/default.qtpl:300
			qw422016.N().S(`
                                <div class="col text-center p-0">
                                `)
//line pkg/report/templates/default.qtpl:302
		}
//line pkg/report/templates/default.qtpl:302
		qw422016.N().S(`
                                    <p class="mx-auto mb-1">`)
//line pkg/report/templates/default.qtpl:303
		qw422016.N().D(sumFail)
//line pkg/report/templates/default.qtpl:303
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
//line pkg/report/templates/default.qtpl:317
		qw422016.E().S(p.ConfigAuditReport.Report.UpdateTimestamp.String())
//line pkg/report/templates/default.qtpl:317
		qw422016.N().S(`</p>
                                </div>
                             </div>
                        </div>

                    </div>      
                </div>
                  <div class="row"><h5 class="text-info" id="ca_pod_checks">Pod Checks</h5></div>
                  <div class="row">
                      <table class="table table-sm table-bordered">
                          <thead>
                              <tr>
                                <th scope="col">PASS</th>
                                <th scope="col">ID</th>
                                <th scope="col">Severity</th>
                                <th scope="col">Category</th>
                              </tr>
                            </thead>
                            <tbody>
                              `)
//line pkg/report/templates/default.qtpl:336
		for _, check := range p.ConfigAuditReport.Report.PodChecks {
//line pkg/report/templates/default.qtpl:336
			qw422016.N().S(`
                                <tr>
                                  <td>`)
//line pkg/report/templates/default.qtpl:338
			qw422016.E().V(check.Success)
//line pkg/report/templates/default.qtpl:338
			qw422016.N().S(`</td>
                                  <td>`)
//line pkg/report/templates/default.qtpl:339
			qw422016.E().S(check.ID)
//line pkg/report/templates/default.qtpl:339
			qw422016.N().S(`</td>
                                  <td>`)
//line pkg/report/templates/default.qtpl:340
			qw422016.E().S(check.Severity)
//line pkg/report/templates/default.qtpl:340
			qw422016.N().S(`</td>
                                  <td>`)
//line pkg/report/templates/default.qtpl:341
			qw422016.E().S(check.Category)
//line pkg/report/templates/default.qtpl:341
			qw422016.N().S(`</td>
                                </tr>
                              `)
//line pkg/report/templates/default.qtpl:343
		}
//line pkg/report/templates/default.qtpl:343
		qw422016.N().S(`
                            </tbody>
                      </table>
                  </div>
                  `)
//line pkg/report/templates/default.qtpl:347
		for container, checks := range p.ConfigAuditReport.Report.ContainerChecks {
//line pkg/report/templates/default.qtpl:347
			qw422016.N().S(`
                    <div class="row"><h5 class="text-info" id="ca_container_`)
//line pkg/report/templates/default.qtpl:348
			qw422016.E().S(container)
//line pkg/report/templates/default.qtpl:348
			qw422016.N().S(`">Container `)
//line pkg/report/templates/default.qtpl:348
			qw422016.E().S(container)
//line pkg/report/templates/default.qtpl:348
			qw422016.N().S(`</h5></div>
                    <div class="row">
                        <table class="table table-sm table-bordered">
                            <thead>
                                <tr>
                                  <th scope="col">PASS</th>
                                  <th scope="col">ID</th>
                                  <th scope="col">Severity</th>
                                  <th scope="col">Category</th>
                                </tr>
                              </thead>
                              <tbody>
                                `)
//line pkg/report/templates/default.qtpl:360
			for _, check := range checks {
//line pkg/report/templates/default.qtpl:360
				qw422016.N().S(`
                                  <tr>
                                    <td>`)
//line pkg/report/templates/default.qtpl:362
				qw422016.E().V(check.Success)
//line pkg/report/templates/default.qtpl:362
				qw422016.N().S(`</td>
                                    <td>`)
//line pkg/report/templates/default.qtpl:363
				qw422016.E().S(check.ID)
//line pkg/report/templates/default.qtpl:363
				qw422016.N().S(`</td>
                                    <td>`)
//line pkg/report/templates/default.qtpl:364
				qw422016.E().S(check.Severity)
//line pkg/report/templates/default.qtpl:364
				qw422016.N().S(`</td>
                                    <td>`)
//line pkg/report/templates/default.qtpl:365
				qw422016.E().S(check.Category)
//line pkg/report/templates/default.qtpl:365
				qw422016.N().S(`</td>
                                  </tr>
                                `)
//line pkg/report/templates/default.qtpl:367
			}
//line pkg/report/templates/default.qtpl:367
			qw422016.N().S(`
                              </tbody>
                        </table>
                    </div>
                  `)
//line pkg/report/templates/default.qtpl:371
		}
//line pkg/report/templates/default.qtpl:371
		qw422016.N().S(`
                  `)
//line pkg/report/templates/default.qtpl:372
	}
//line pkg/report/templates/default.qtpl:372
	qw422016.N().S(`
            </div>
        </div>
`)
//line pkg/report/templates/default.qtpl:375
}

//line pkg/report/templates/default.qtpl:375
func (p *ReportPage) WriteBody(qq422016 qtio422016.Writer) {
//line pkg/report/templates/default.qtpl:375
	qw422016 := qt422016.AcquireWriter(qq422016)
//line pkg/report/templates/default.qtpl:375
	p.StreamBody(qw422016)
//line pkg/report/templates/default.qtpl:375
	qt422016.ReleaseWriter(qw422016)
//line pkg/report/templates/default.qtpl:375
}

//line pkg/report/templates/default.qtpl:375
func (p *ReportPage) Body() string {
//line pkg/report/templates/default.qtpl:375
	qb422016 := qt422016.AcquireByteBuffer()
//line pkg/report/templates/default.qtpl:375
	p.WriteBody(qb422016)
//line pkg/report/templates/default.qtpl:375
	qs422016 := string(qb422016.B)
//line pkg/report/templates/default.qtpl:375
	qt422016.ReleaseByteBuffer(qb422016)
//line pkg/report/templates/default.qtpl:375
	return qs422016
//line pkg/report/templates/default.qtpl:375
}
