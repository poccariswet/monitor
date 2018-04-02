yum install -y vim
yum install -y git
yum install -y nginx
wget https://dl.google.com/go/go1.10.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz

# echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
# echo "export GOPATH=$HOME/go" >> ~/.bashrc
# echo "export PATH=$PATH:$GOPATH/bin" >> ~/.bashrc
