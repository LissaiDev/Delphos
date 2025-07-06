export interface Host {
    hostname: string;
    os: string;
    uptime: number;
  }
  
  export interface Memory {
    total: number;
    used: number;
    free: number;
    swapTotal: number;
    swapUsed: number;
    swapFree: number;
  }
  
  export interface CPU {
    usage: number;
    model: string;
    cores: number;
  }
  
  export interface Disk {
    mountpoint: string;
    type: string;
    total: number;
    used: number;
    free: number;
    usedPercent: number;
  }
  
  export interface Network {
    totalBytesSent: number;
    totalBytesRecv: number;
    interfaceName: string;
  }
  
  export interface Monitor {
    host: Host;
    memory: Memory;
    cpu: CPU[];
    disk: Disk[];
    network: Network[];
  }
  