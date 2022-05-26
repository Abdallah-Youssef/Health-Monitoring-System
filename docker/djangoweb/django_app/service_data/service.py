class service:
    def __init__(self,service_name, ram_total,cpu,ram_free,disk_total,disk_free,count,peak_cpu,peak_ram_free,
                 peak_ram_total,peak_disk_total,peak_disk_free):
        self.service_name = service_name
        self.ram_total = ram_total
        self.cpu= cpu
        self.ram_free = ram_free
        self.disk_total = disk_total
        self.disk_free = disk_free
        self.count = count
        self.peak_cpu=peak_cpu
        self.peak_ram_free = peak_ram_free
        self.peak_ram_total = peak_ram_total
        self.peak_disk_total = peak_disk_total
        self.peak_disk_free = peak_disk_free
