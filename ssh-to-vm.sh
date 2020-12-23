#!/bin/bash

ssh -R 52698:localhost:52698 -i azure-vm.pem hlf@23.101.222.137
