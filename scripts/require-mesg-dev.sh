if [ ! "$MESG_DEV" = "true" ]; then
   echo "you must run scripts via ./mesg-dev script,";
   echo "try executing the following command from the root of core:";
   echo -e "\t./mesg-dev $0 $*"
   exit 1;
fi