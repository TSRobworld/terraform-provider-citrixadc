/*
* Copyright (c) 2021 Citrix Systems, Inc.
*
*   Licensed under the Apache License, Version 2.0 (the "License");
*   you may not use this file except in compliance with the License.
*   You may obtain a copy of the License at
*
*       http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*   distributed under the License is distributed on an "AS IS" BASIS,
*   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*   See the License for the specific language governing permissions and
*   limitations under the License.
*/

package ns

/**
* Configuration for "VPX" resource.
*/
type Nsvpxparam struct {
	/**
	* This setting applicable in virtual appliances, to move master clock source cpu from management cpu cpu0 to cpu1 ie PE0.
		* There are 2 options for the behavior:
		1. YES - Allow the Virtual Appliance to move clock source to cpu1.
		2. NO - Virtual Appliance will use management cpu ie cpu0 for clock source default option is NO.
	*/
	Masterclockcpu1 string `json:"masterclockcpu1,omitempty"`
	/**
	* This setting applicable in virtual appliances, is to affect the cpu yield(relinquishing the cpu resources) in any hypervised environment.
		* There are 3 options for the behavior:
		1. YES - Allow the Virtual Appliance to yield its vCPUs periodically, if there is no data traffic.
		2. NO - Virtual Appliance will not yield the vCPU.
		3. DEFAULT - Restores the default behaviour, according to the license.
		* Its behavior in different scenarios:
		1. As this setting is node specific only, it will not be propagated to other nodes, when executed on Cluster(CLIP) and HA(Primary).
		2. In cluster setup, use '-ownerNode' to specify ID of the cluster node.
		3. This setting is a system wide implementation and not granular to vCPUs.
		4. No effect on the management PE.
	*/
	Cpuyield string `json:"cpuyield,omitempty"`
	/**
	* ID of the cluster node for which you are setting the cpuyield. It can be configured only through the cluster IP address.
	*/
	Ownernode uint32 `json:"ownernode,omitempty"`

	//------- Read only Parameter ---------;

	Vpxenvironment string `json:"vpxenvironment,omitempty"`
	Memorystatus string `json:"memorystatus,omitempty"`
	Cloudproductcode string `json:"cloudproductcode,omitempty"`
	Vpxoemcode string `json:"vpxoemcode,omitempty"`
	Technicalsupportpin string `json:"technicalsupportpin,omitempty"`

}
