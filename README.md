# Discovery script for DNS domains

Outputs information about a DNS domain.

This project is designed to output a high-level overview of a DNS zone. For instance, given `example.com` it should output 

    - A quick overview of the site with an executive summary of the information listed below. 
    - List the zone name, registrar, and nameservers
      - Try to determine the "friendly" name for the registrar and nameservers. For example, domaincontrol.com is run by GoDaddy, so if the Nameserver records use domaincontrol, say that the site is using GoDaddy for DNS management.
    - A list of configured services, hosts, and redirects
      - For each configured service, it should do some sort of quick health check
        - As an example, if email is configured, it should do a quick health check to ensure email DNS entries are in place and valid for mail, DKIM, and SPF. 
        - For A and CNAME records, it should check if the target is publicly facing, has a valid certificate, and is using TLS version 1.2+

# User interaction

I'm envisioning this as a CLI-first tool, although it could be loaded as a libary. Python has a readily-available pattern for doing this, and I'll use the click library to help make the CLI usage easy.

# Output format

For now, use Markdown, and output to a directory called `output`, in the project root location. Eventually I want a really nice, HTML-based output. I'm not a UI guy, so this will be last, and I'll probably be lazy about it. Actually, I'll probably use mk-docs at first. 

