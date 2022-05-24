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
    formated_str = "<button><a href=""/"">make another query</a></button><div id=""duration"">"+"<h1>"+str(start_date)+" to "+str(end_date)+"</h1>"+"</div><br>"
    for row in services_data:
        count = row.count
        formated_str = formated_str+"<h1>" +str(row.service_name)+"</h1>"+\
                       "Peak CPU:"+str(row.cpu)+"<br>"+\
                       "Avg CPU:"+str(row.cpu/count)+"<br>"+\
                       "Avg Ram total:"+str(row.ram_total/count)+"<br>"+\
                       "Peak Ram total:"+str(row.ram_total)+"<br>"+\
                       "Avg Ram free:"+str(row.ram_free/count)+"<br>"+\
                       "Peek Ram free:"+str(row.ram_free)+"<br>"+\
                       "Avg Disk total:"+str(row.disk_total/count)+"<br>"+\
                       "Peek Disk total:"+str(row.disk_total)+"<br>"+\
                       "Avg Disk free:"+str(row.disk_free/count)+"<br>"+\
                       "Peak Disk free:"+str(row.disk_free)+"<br>"
    return formated_str


def format_dates(start_date,end_date)->list:
    start_date_splitted = start_date.split()
    start_date = start_date_splitted[0]+" "+start_date_splitted[1]
    end_date = start_date_splitted[2]+" "+end_date
    return [start_date, end_date]
