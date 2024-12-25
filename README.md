# Advanced Port Scanner 🔍

A modern and user-friendly web-based port scanning tool. This application is a professional port scanner developed for network security testing and system administration.


## 🚀 Features

- **Modern Web Interface**: User-friendly, responsive design
- **Multiple Port Scanning Options**:
  - Common Ports (20-23, 25, 53, 80, etc.)
  - Well-Known Ports (1-1024)
  - Custom Port Range
  - Specific Ports
- **Advanced Features**:
  - Concurrent Scanning
  - Timeout Management
  - Detailed Logging
  - Result Copying
- **Security Features**:
  - IP Address Validation
  - Port Range Validation
  - Error Handling

## 🛠️ Requirements

- Python 3.8+
- Go 1.16+
- Modern web browser (Chrome, Firefox, Safari, Edge)

## ⚙️ Installation

1. Clone the repository:
```bash
git clone https://github.com/yusufdalbudak/advanced-port-scanner.git
cd port-scanner
```

2. Install required Python packages:
```bash
pip install -r requirements.txt
```

3. Start the application:
```bash
go run main.go
```

4. Open in your web browser:
```
http://localhost:8080
```

## 🎯 Usage

1. Enter IP address (e.g., 192.168.1.1)
2. Select port scanning option:
   - Common Ports
   - Well-Known Ports (1-1024)
   - Custom Port Range
   - Specific Ports
3. Click "Start Scan"
4. View results and copy if needed

## 📝 Example Use Cases

### Scanning Common Ports
```
IP: 192.168.1.1
Option: Common Ports
```

### Scanning Custom Port Range
```
IP: 192.168.1.1
Option: Custom Range
Port Range: 80-443
```

### Scanning Specific Ports
```
IP: 192.168.1.1
Option: Specific Ports
Ports: 80,443,3306,8080
```

## 🔒 Security Notes

- Use this tool only on networks where you have permission
- Port scanning may be prohibited on some networks
- Firewalls may block scanning attempts
- High volume port scanning may slow down the network

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- Yusuf Dalbudak - [GitHub](https://github.com/yusufdalbudak)

## 🙏 Acknowledgments

- [Go](https://golang.org/) - Main programming language
- [Bootstrap](https://getbootstrap.com/) - UI framework
- [Font Awesome](https://fontawesome.com/) - Icons

## 📞 Contact

For questions, please use the Issues section or reach out through:

- Email: yusufdalbudak2121@gmail.com
- Website: [cybersecdev.com](https://cybersecdev.com)

## 📊 Project Status

![GitHub stars](https://img.shields.io/github/stars/yusufdalbudak/port-scanner?style=social)
![GitHub forks](https://img.shields.io/github/forks/yusufdalbudak/port-scanner?style=social)
![GitHub issues](https://img.shields.io/github/issues/yusufdalbudak/port-scanner)
![GitHub license](https://img.shields.io/github/license/yusufdalbudak/port-scanner)
