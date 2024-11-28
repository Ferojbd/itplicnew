# Automatically generated Red Hat Update Agent config file, do not edit.
# Format: 1.0
enableProxy[comment]=Use a HTTP Proxy
enableProxy=0

serverURL[comment]=Remote server URL (use FQDN)
serverURL=http://space.itplic.biz/XMLRPC/

mirrorURL[comment]=Mirror list URL
mirrorURL=https://api.itplic.biz/cln/mirror

debug[comment]=Whether or not debugging is enabled
debug=0

systemIdPath[comment]=Location of system id
systemIdPath=/etc/sysconfig/rhn/systemid

versionOverride[comment]=Override the automatically determined system version
versionOverride=

httpProxy[comment]=HTTP proxy in host:port format, e.g. squid.redhat.com:3128
httpProxy=:0

proxyUser[comment]=The username for an authenticated proxy
proxyUser=

proxyPassword[comment]=The password to use for an authenticated proxy
proxyPassword=

enableProxyAuth[comment]=To use an authenticated proxy or not
enableProxyAuth=0

networkRetries[comment]=Number of attempts to make at network connections before giving up
networkRetries=1

sslCACert[comment]=The CA cert used to verify the ssl server
sslCACert=/usr/share/rhn/CLN-ORG-TRUSTED-SSL-CERT

noReboot[comment]=Disable the reboot actions
noReboot=0

disallowConfChanges[comment]=Config options that can not be overwritten by a config update action
disallowConfChanges=noReboot;sslCACert;useNoSSLForPackages;noSSLServerURL;serverURL;disallowConfChanges

retrieveOnly[comment]=Retrieve packages only
retrieveOnly=0

writeChangesToLog[comment]=Log to /var/log/up2date which packages has been added and removed
writeChangesToLog=0

stagingContentWindow[comment]=How much forward we should look for future actions. In hours.
stagingContentWindow=24

useNoSSLForPackages[comment]=Use the noSSLServerURL for package, package list, and header fetching (disable Akamai)
useNoSSLForPackages=0

tmpDir[comment]=Use this Directory to place the temporary transport files
tmpDir=/tmp

skipNetwork[comment]=Skips network information in hardware profile sync during registration.
skipNetwork=0

stagingContent[comment]=Retrieve content of future actions in advance
stagingContent=1

hostedWhitelist[comment]=RHN Hosted URL's
hostedWhitelist=

clnServerURL[comment]=None
clnServerURL=http://cln.cloudlinux.com/clweb/xmlrpc

