#!/bin/bash


REGIONS="eu-west-1 eu-north-1 us-east-2"
ACCOUNTS="012345678910 012345678911"
NAMES="front front front db db db redis service webapp nginx proxy worker"
STATES="stopped running running running"
TYPES="t2.micro t4g.nano m6g.xlarge c5a.large"

# static random vmlist html
echo '<div class="container-fluid">
<table class="table table-hover" id="vmlist-table">
<thead>
  <tr>
    <th></th>
    <th>Region</th>
    <th>AwsAccountId</th>
    <th>InstanceId</th>
	  <th>Name</th>
	  <th>Architecture</th>
	  <th>PlatformDetails</th> 
	  <th>InstanceType</th>
	  <th>State</th>       
  </tr>
</thead>
<tbody>' > vmlist

for num in {001..100}
do
    ACCOUNT=$(shuf -e -n1 ${ACCOUNTS})
    REGION=$(shuf -e -n1 ${REGIONS})
    TYPE=$(shuf -e -n1 ${TYPES})
    NAME="$(shuf -e -n1 ${NAMES})-${num}"
    ID="i-0d2941b87e51def${num}"
    STATE=$(shuf -e -n1 ${STATES})
    if [ $TYPE == "t2.micro" ] || [ $TYPE == "c5a.large" ]
    then
        ARCH="x86_64"
    else
        ARCH="arm_64"
    fi

    echo "<tr>
        <td>
          <div class=\"form-check\">
          <input class=\"form-check-input\" type=\"radio\" name=\"flexRadio\" id=\"flexRadioDefault${ID}\"
           hx-get=\"/cloud-commis/vmdetails-${ARCH}\"
           hx-trigger=\"click\" hx-target=\"#vmdetails\" hx-swap=\"innerHTML\">
         </div>
       </td>
       <td>${REGION}</td>
       <td>${ACCOUNT}</td>
       <td>${ID}</td>
	   <td>${NAME}</td>
	   <td>${ARCH}</td>
	   <td>Linux/UNIX</td>  
	   <td>${TYPE}</td>
	   <td>${STATE}</td>
       </tr>" >> vmlist
done



echo "</tbody>
</table>
</div>
<script>
\$('#vmlist-table').DataTable();
</script>
<div id=vmdetails class=\"border border-secondary text-bg-light d-flex dropdown\" >
</div>" >> vmlist


# static  vmdetails-arm_64 html
echo '<div class="container p-2 d-flex justify-content-center">
<div class="row overflow-auto gy-5">

    <div class="col"> 
Instance details :
    <pre class="p-3 col border bg-opacity-10 bg-warning" style="width:600px; height: 400px;">

{
  "Name": "test",
  "Architecture": "arm64",
  "LaunchTime": "2024-10-26T13:15:55Z",
  "UsageOperationUpdateTime": "2024-10-26T13:15:55Z",
  "PlatformDetails": "Linux/UNIX",
  "ImageId": "ami-01c5300f289d64643",
  "InstanceType": "t4g.nano",
  "PublicIpAddress": "",
  "State": "stopped",
  "KernelImage": "Kernel 6.1.112-122.189.amzn2023.aarch64 on an aarch64 (-)\r",
  "BootImage": "Booting `Amazon Linux (6.1.112-122.189.amzn2023.aarch64) 2023"
}

    </pre>
    </div>

    <div class="col">
Image details :
    <pre class="p-3 col border bg-opacity-10 bg-warning"" style="width:600px; height: 400px;">

{
  "Name": "al2023-ami-2023.6.20241010.0-kernel-6.1-arm64",
  "Region": "eu-west-3",
  "Description": "Amazon Linux 2023 AMI 2023.6.20241010.0 arm64 HVM kernel-6.1",
  "OwnerId": "137112412989",
  "DeprecationTime": "2025-01-09T16:54:00.000Z"
}

    </pre>
    </div>

  </div>
</div>' > vmdetails-arm_64

# static  vmdetails-x86_64 html
echo '
<div class="container p-2 d-flex justify-content-center">
  <div class="row overflow-auto gy-5">

    <div class="col"> 
Instance details :
    <pre class="p-3 col border bg-opacity-10 bg-warning" style="width:600px; height: 400px;">

{
  "Name": "test2",
  "Architecture": "x86_64",
  "LaunchTime": "2024-10-26T13:31:35Z",
  "UsageOperationUpdateTime": "2024-10-26T13:31:35Z",
  "PlatformDetails": "Linux/UNIX",
  "ImageId": "ami-07db896e164bc4476",
  "InstanceType": "t2.micro",
  "PublicIpAddress": "",
  "State": "stopped",
  "KernelImage": "vmlinuz-6.8.0-1015-aws root",
  "BootImage": "Welcome to \u001b[1mUbuntu 22.04.5 LTS\u001b[0m!"
}

    </pre>
    </div>

    <div class="col">
Image details :
    <pre class="p-3 col border bg-opacity-10 bg-warning"" style="width:600px; height: 400px;">

{
  "Name": "ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-20240927",
  "Region": "eu-west-3",
  "Description": "Canonical, Ubuntu, 22.04 LTS, amd64 jammy image build on 2024-09-27",
  "OwnerId": "099720109477",
  "DeprecationTime": "2026-09-27T03:11:33.000Z"
}

    </pre>
    </div>

  </div>
</div>' > vmdetails-x86_64