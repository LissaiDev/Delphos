interface Host {
    hostname: string;
    os: string;
    uptime: number;
  }
  
  interface Memory {
    total: number;
    used: number;
    free: number;
    swapTotal: number;
    swapUsed: number;
    swapFree: number;
  }
  
  interface CPU {
    usage: number;
    model: string;
    cores: number;
  }
  
  interface Disk {
    mountpoint: string;
    type: string;
    total: number;
    used: number;
    free: number;
    usedPercent: number;
  }
  
  interface Network {
    totalBytesSent: number;
    totalBytesRecv: number;
    interfaceName: string;
  }
  
  interface Monitor {
    host: Host;
    memory: Memory;
    cpu: CPU[];
    disk: Disk[];
    network: Network[];
  }
  