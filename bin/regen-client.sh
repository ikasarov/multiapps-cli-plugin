#!/bin/bash

GREEN='\033[0;32m'
BORING='\033[0m' # No Color

check_dir() {
	if [ -d "$1" ]; then
		echo -e "Folder $1 already exists.\nDo you want to ${GREEN}overwrite${BORING} the files?"
		
		read answer
		
		if [ "$answer" != 'y' ] && [ "$answer" != 'Y' ] && [ "$answer" != 'yes' ] && [ "$answer" != 'YES' ]; then
			echo -e "Edit regen folder name and ${GREEN}restart${BORING} script."
			exit 1
		fi
		
		rm -r "$1"
	fi
	
	echo -e "${GREEN}Recreating folder...${BORING}"
	mkdir "$1"
}

go_home=$(echo ~/go)
git_home=$(echo ~/git)
tmp="swagger-regen"
swagger_file="mtarest.yaml"
swagger_file_v2="mtarest_v2.yaml"
client_name="mtaclient"

if [ $# -eq 0 ];
then
    echo "No arguments supplied, generating in dirs ${tmp} and ${tmp}-v2"
elif [ "$1" ]; then
	tmp="$1"
fi

regen_folder="${go_home}/src/github.com/cloudfoundry-incubator/multiapps-cli-plugin/${tmp}"
regen_folder_v2="${regen_folder}-v2"
definition_file="${git_home}/multiapps-controller/com.sap.cloud.lm.sl.cf.api/src/main/resources/${swagger_file}"
definition_file_v2="${git_home}/multiapps-controller/com.sap.cloud.lm.sl.cf.api/src/main/resources/${swagger_file_v2}"

check_dir "${regen_folder}"
check_dir "${regen_folder_v2}"

echo -e "Assuming controller project is under this parent dir: ${GREEN}${git_home}${BORING}"
echo -e "Assuming plugin project is under this parent dir: ${GREEN}${go_home}${BORING}"
echo -e "Reading from\n\t${GREEN}${definition_file}\n\t${definition_file_v2}${BORING}\nGenerating in \n\t${GREEN}${regen_folder}\n\t${regen_folder_v2}${BORING}"

cd "${git_home%/*}"
mvn -f "${git_home}/multiapps-controller/" clean package -DskipTests=true -pl=com.sap.cloud.lm.sl.cf.api

swagger generate client -f ${definition_file} -A http_mta -c ${client_name} -t ${regen_folder}
swagger generate client -f ${definition_file_v2} -A http_mta_v2 -c "${client_name}_v2" -t ${regen_folder_v2}