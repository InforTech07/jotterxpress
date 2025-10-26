%global goipath %{gopath}/src/%{import_path}
%global debug_package %{nil}

Name:           jotterxpress
Version:        1.0.0
Release:        1%{?dist}
Summary:        Fast and simple CLI tool for taking notes
License:        MIT
URL:            https://github.com/yourusername/jotterxpress
Source0:        %{name}-%{version}.tar.gz

# BuildRequires:  golang >= 1.20  # Commented out if Go is installed manually
# BuildRequires:  git

%description
JotterXpress is a fast and simple CLI tool for taking notes, managing tasks, 
contacts, and reminders. It provides an interactive terminal interface with 
beautiful styling and efficient note management.

%prep
%setup -q -n %{name}-%{version}

%build
%gobuild -o jotterxpress cmd/jotterxpress/main.go

%install
mkdir -p %{buildroot}%{_bindir}
install -m 0755 jotterxpress %{buildroot}%{_bindir}/jtx

# Create directory for notes in user's home
mkdir -p %{buildroot}%{_sysconfdir}/skel/.jotterxpress/notes

%files
%defattr(-,root,root,-)
%{_bindir}/jtx
%config(noreplace) %{_sysconfdir}/skel/.jotterxpress/notes

%clean
rm -rf %{buildroot}

%changelog
* Sat Oct 25 2025 Your Name <your.email@example.com> - 1.0.0-1
- Initial release of JotterXpress
