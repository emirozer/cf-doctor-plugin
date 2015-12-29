#CloudFoundry CLI Plugin - Doctor

This plugin is obviously inspired from [brew](http://brew.sh/) doctor :) it will scan your cloudfoundry to see if there are anomalies or useful action points that it can report back to you. Current functionality is only focused on apps and routes..
This plugin does *not* change any state or configuration, it merely just scans and gathers information than reports back anomalies.

List of all plugins available: <http://plugins.cloudfoundry.org/ui/>

### Installation

Get the latest release in binaries depending on your os/arch here: <https://github.com/emirozer/cf-doctor-plugin/releases>

Run **cf install-plugin BINARY_FILENAME** to install a plugin. Replace **BINARY_FILENAME** with the path to and name of the binary file.

After installation just run:

    cf doctor


<br>
##### Sample output
![Screenshot](https://raw.github.com/emirozer/cf-doctor-plugin/master/docs/ndoc.png)
<br>

##Licence

 Copyright 2015 Emir Ozer

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
