import datetime


def format_result(services_data,start_date,end_date):
    """ format the query result into HTML string to be rendered
    Paramaters
    ------------
    services_data:
            the resulted query of services
    start_date:
            the start_date to query about
        end_date:
            the end_date to query about
    Retuens
    ----------
    String
        a formatted HTML string of the services query data
    """
    #todo format the result of the query
    styling="<html><style>table, th, td {border:1px solid black;}</style><body>"
    formated_str = styling+ "<button><a href=""/"">make another query</a></button><div id=""duration"">"+"<h1>"+str(start_date)+" to "+str(end_date)+"</h1>"+"</div><br>"
    for row in services_data:
        count = row.count
        formated_str = formated_str+"<h1>" +str(row.service_name)+"</h1>"+\
                       "<table><tr><td>Peak CPU:</td><td>"+str(row.cpu)+"</td></tr>"+\
                       "<tr><td>Avg CPU:</td><td>"+str(row.cpu/count)+"</td></tr>"+\
                       "<tr><td>Avg Ram total:</td><td>"+str(row.ram_total/count)+"</td></tr>"+\
                       "<tr><td>Peak Ram total:</td><td>"+str(row.ram_total)+"</td></tr>"+\
                       "<tr><td>Avg Ram free:</td><td>"+str(row.ram_free/count)+"</td></tr>"+\
                       "<tr><td>Peek Ram free:</td><td>"+str(row.ram_free)+"</td></tr>"+\
                       "<tr><td>Avg Disk total:</td><td>"+str(row.disk_total/count)+"</td></tr>"+\
                       "<tr><td>Peek Disk total:</td><td>"+str(row.disk_total)+"</td></tr>"+\
                       "<tr><td>Avg Disk free:</td><td>"+str(row.disk_free/count)+"</td></tr>"+\
                       "<tr><td>Peak Disk free:</td><td>"+str(row.disk_free)+"</td></tr></table></body></html>"
    return formated_str


def format_dates(start_date,end_date)->list:
    start_date_splitted = start_date.split()
    s_days,s_month,s_year=start_date_splitted[0].split("/")
    s_hours,s_minutes=start_date_splitted[1].split(":")
    e_days,e_month,e_year = start_date_splitted[2].split("/")
    e_hours, e_minutes = end_date.split(":")
    start_date = datetime.datetime(int(s_year),int(s_month),int(s_days),int(s_hours),int(s_minutes))
    end_date = datetime.datetime(int(e_year),int(e_month),int(e_days),int(e_hours),int(e_minutes))
    return [start_date, end_date]
