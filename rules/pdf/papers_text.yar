rule Paper : PDF
{
    meta:
        version = "0.1"
        description = "generic rule for detecting papers"
    
    strings:
        $abstract_keyword = "Abstract" nocase
        $introduction_keyword = "Introduction" nocase
        $conclusion_keyword = "Conclusion" nocase
        
        $experiment_keyword = "Experiment" nocase
        $evaluation_keyword = "Evaluation" nocase
        

    condition:
        $abstract_keyword and any of ($introduction_keyword, $conclusion_keyword, $experiment_keyword, $evaluation_keyword)
}

rule NDSSPaper : PDF
{
    meta:
        version = "0.1"

    strings:
        $ndss_keyword = "Network and Distributed Systems Security (NDSS) Symposium"
        $ndss_keyword2 = "Network and Distributed System Security (NDSS) Symposium"
        
        $url = "www.ndss-symposium.org"
    condition:
        ($ndss_keyword or $ndss_keyword2) 
        and $url
}

rule USENIXSecurityPaper : PDF
{
    meta:
        version = "0.1"

    strings:
        $sponsor_keyword = "sponsored by USENIX"
        $sec = "USENIX Security Symposium"
        $raid = "International Symposium on Research in Attacks, Intrusions and Defenses"

    condition:
        $sponsor_keyword and ($sec or $raid)
}

rule PreprintPaper : PDF
{
    meta:
        version = "0.1"

    strings:
        $arxiv_pattern = /arXiv:\d.*\[.*\].* \d{4}/

    condition:
        $arxiv_pattern
}

rule ACMPaper : PDF
{
    meta:
    version = "0.1"

    strings:
        $acm_keyword = "ACM Reference Format"

    condition:
        $acm_keyword
}

rule IEEEPaper : PDF
{
    meta:
        version = "0.1"

    strings:
        $ieee_copyright_keyword = "Personal use is permitted, but republication/redistribution requires IEEE permission."
        $conference_keyword = /^\d{4} IEEE.*/ nocase
        $keywords = "IEEE TRANSACTIONS ON SOFTWARE ENGINEERING, VOL"

    condition:
        $ieee_copyright_keyword or $conference_keyword or $keywords
}

rule CCSPaper
{
    meta:
        version = "0.1"

    strings:
        $magic = { 25 50 44 46 }
        //$ccs_keyword = "CCS"
        $ccs_logo = "ccs_logo"

    condition:
        $magic at 0 and any of ($ccs_logo)
}

rule NaturePaper
{
    meta:
        version = "0.1"

    strings:
        $magic = { 25 50 44 46 }
        $nature_reprint = "http://www.nature.com/reprints"
        $doi_keyword = "https://doi.org/"

    condition:
        $magic at 0 and ($nature_reprint and $doi_keyword)
}