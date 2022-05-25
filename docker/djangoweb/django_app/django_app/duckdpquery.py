import duckdb
from service_data.service import service

def add_service_data(main_service,instance):
    if main_service.service_name == "":
        main_service.service_name = instance[1]
    main_service.ram_total = main_service.ram_total+instance[2]
    main_service.cpu = main_service.cpu+instance[3]
    main_service.ram_free = main_service.ram_free+instance[4]
    main_service.disk_total = main_service.disk_total+instance[5]
    main_service.disk_free = main_service.disk_free+instance[6]
    main_service.count = main_service.count+instance[7]
    return main_service

def query(start_date,end_date):
    cursor = duckdb.connect()
    services = cursor.execute(
    "SELECT * FROM 'django_app/output/*.parquet' Where Timestamp BETWEEN (?) AND (?) ",
    [start_date, end_date]).fetchall()
    query_result=[]
    for i in range(4):
        query_result.append(service("", 0,0,0,0,0,0))
    for instance in services:
        name=instance[1]
        match name:
            case "service-1":
                query_result[0]=add_service_data(query_result[0],instance)
            case "service-2":
                query_result[1] = add_service_data(query_result[1], instance)
            case "service-3":
                query_result[2] = add_service_data(query_result[2], instance)
            case "service-4":
                query_result[3] = add_service_data(query_result[3], instance)
    query_filtered=[]
    for instance in query_result:
        if not instance.service_name == "":
            query_filtered.append(instance)

    return query_filtered
