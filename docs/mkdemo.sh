#!/bin/bash


REGIONS="eu-west-1 eu-north-1 us-east-2"
ACCOUNTS="012345678910 012345678911"
NAMES="front front front db db db redis service webapp nginx proxy worker"
STATES="stopped running running running"
TYPES="t2.micro t4g.nano m6g.xlarge c5a.large"

# static random vmlist html
echo '<div class="container-fluid">
<table class="table table-hover" id="vmlist-table" style="font-size: 12px;">
<thead>
  <tr>
    <th>Region</th>
    <th>AwsAccountId</th>
    <th>InstanceId</th>
	  <th>Name</th>
	  <th>Architecture</th>
	  <th>PlatformDetails</th> 
	  <th>InstanceType</th>
	  <th>State</th>
    <th>Details</th>   
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
       <td>${REGION}</td>
       <td>${ACCOUNT}</td>
       <td>${ID}</td>
	   <td>${NAME}</td>
	   <td>${ARCH}</td>
	   <td>Linux/UNIX</td>  
	   <td>${TYPE}</td>
	   <td>${STATE}</td>
     <td>
       <button class="btn btn-primary" type="button" data-bs-toggle="offcanvas" data-bs-target="#offcanvasRight" aria-controls="offcanvasRight"
         hx-get="/cloud-commis/vmdetails"
         hx-trigger="click" hx-target="#rightPanelDetails" hx-swap="innerHTML">
         Details
       </button>
     </td>
       </tr>" >> vmlist
done



echo "</tbody>
</table>
</div>
<script>
\$('#vmlist-table').DataTable();
</script>
<div id=vmdetails class=\"border border-secondary text-bg-light d-flex dropdown\" >
</div>">> vmlist

echo '<div class="offcanvas offcanvas-end" data-bs-scroll="true" 
data-bs-backdrop="false" tabindex="-1" 
id="offcanvasRight" aria-labelledby="offcanvasRightLabel" >
  <div class="offcanvas-header">
    <h5 id="offcanvasRightLabel">vm Details</h5>
    <button type="button" class="btn-close text-reset" data-bs-dismiss="offcanvas" aria-label="Close"></button>
  </div>
  <div class="offcanvas-body">
    <!-- style="overflow:auto; width: 350px;" -->
    <pre class="bg-opacity-10 bg-warning " id="rightPanelDetails" style="white-space: pre-wrap;">
    </pre>
  </div>
</div>'>> vmlist


# static  vmdetails-arm_64 html
echo '
{
  "Name": "demo values",
  "Architecture": "arm64",
  "LaunchTime": "2024-10-26T13:15:55Z",
  "UsageOperationUpdateTime": "2024-10-26T13:15:55Z",
  "PlatformDetails": "Linux/UNIX",
  "ImageId": "ami-01c5300f289d6XXXX",
  "InstanceType": "xxx.someSize",
  "PublicIpAddress": "",
  "State": "running",
  "KernelImage": "Kernel 6.1.112-122.189.amzn2023.aarch64 on an aarch64 (-)\r",
  "BootImage": "Booting `Amazon Linux (6.1.112-122.189.amzn2023.aarch64) 2023"
}

{
  "Name": "al2023-ami-2023.6.20241010.0-kernel-6.1-arm64",
  "Region": "xx-west-3",
  "Description": "Amazon Linux 2023 AMI 2023.6.20241010.0 arm64 HVM kernel-6.1",
  "OwnerId": "137112412989",
  "DeprecationTime": "2025-01-09T16:54:00.000Z"
}
' > vmdetails