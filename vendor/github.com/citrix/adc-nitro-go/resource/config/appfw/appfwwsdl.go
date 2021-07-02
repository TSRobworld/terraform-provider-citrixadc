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

package appfw

/**
* Configuration for WSDL file resource.
*/
type Appfwwsdl struct {
	/**
	* Name of the WSDL file to remove.
	*/
	Name string `json:"name,omitempty"`
	/**
	* URL (protocol, host, path, and name) of the WSDL file to be imported is stored.
		NOTE: The import fails if the object to be imported is on an HTTPS server that requires client certificate authentication for access.
	*/
	Src string `json:"src,omitempty"`
	/**
	* Any comments to preserve information about the WSDL.
	*/
	Comment string `json:"comment,omitempty"`
	/**
	* Overwrite any existing WSDL of the same name.
	*/
	Overwrite bool `json:"overwrite,omitempty"`

	//------- Read only Parameter ---------;

	Response string `json:"response,omitempty"`

}
