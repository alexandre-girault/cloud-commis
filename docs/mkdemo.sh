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
           hx-get=\"/vmdetails/${ACCOUNT}/${REGION}/${ID}\"
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

echo '</tbody>
</table>
</div>
<script>
$('#vmlist-table').DataTable();
</script>

<div id=vmdetails class="border border-secondary text-bg-light d-flex dropdown" >
</div>' >> vmlist