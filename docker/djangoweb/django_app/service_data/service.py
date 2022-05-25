class service:
    def __init__(self,service_name, ram_total,cpu,ram_free,disk_total,disk_free,count):
        self.service_name = service_name
        self.ram_total = ram_total
        self.cpu= cpu
        self.ram_free = ram_free
        self.disk_total = disk_total
        self.disk_free = disk_free
        self.count = count