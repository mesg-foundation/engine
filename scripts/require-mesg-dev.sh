if [ ! "$MESG_DEV" = "true" ]; then
   echo "you must run scripts via ./mesg-dev script,";
   echo "try executing this command from the root of core: ./mesg-dev $0 $*";
   exit 1;
fi