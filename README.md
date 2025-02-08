<a id="readme-top"></a>

<br />
<div align="center">

<h3 align="center">Cloud Commis</h3>

  <p align="center">
    A simple tool to track and audit virtual machines on cloud 
    <br />
    <a href="https://github.com/alexandre-girault/cloud-commis"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/alexandre-girault/cloud-commis">View Demo</a>
    &middot;
    <a href="https://github.com/alexandre-girault/cloud-commis/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    &middot;
    <a href="https://github.com/alexandre-girault/cloud-commis/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

[![Product Name Screen Shot][product-screenshot]](https://example.com)

Here's a blank template to get started.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

* [![Golang]][Golang-logo][Golang-url]
* [![Bootstrap][Bootstrap.com]][Bootstrap-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

Cloud commis could be use locally or in a remote server or container orchestrator.

### Prerequisites

You need a working aws cli. See https://docs.aws.amazon.com/cli/v1/userguide/cli-configure-files.html for configuration.
By default cloud commis will use the current aws credentials. 


### Running from binary 
1. download a tagged release from github
    ```sh
    curl -OL https://github.com/alexandre-girault/cloud-commis/releases/download/v0.0.5/cloudcommis-linux-v0.0.5-amd64
    ```
2. We create a directory to store scan results
    ```sh
    mkdir ~/cloudcommis-data/
    ```
3. Launch cloud-commis
    ```sh
    export CC_localStoragePath="~/cloudcommis-data/" 
    chmod +x cloudcommis-linux-v0.0.5-amd64 
    ./cloudcommis-linux-v0.0.5-amd64 --version
    ./cloudcommis-linux-v0.0.5-amd64

    ```
4. See results on http://localhost:8080/


### Running from source code
1. Clone the repo
   ```sh
   git clone https://github.com/alexandre-girault/cloud-commis.git
   ```
2. run te code (with config.yaml provided at root of this repo)
   ```sh
   make run
   ```
3. See results on http://localhost:8080/


<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources.

_For more examples, please refer to the [Documentation](https://example.com)_

<p align="right">(<a href="#readme-top">back to top</a>)</p>





<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Top contributors:

<a href="https://github.com/alexandre-girault/cloud-commis/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=alexandre-girault/cloud-commis" alt="contrib.rocks image" />
</a>



<!-- LICENSE -->
## License

Distributed under the project_license. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Your Name - [@twitter_handle](https://twitter.com/twitter_handle) - email@email_client.com

Project Link: [https://github.com/alexandre-girault/cloud-commis](https://github.com/alexandre-girault/cloud-commis)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* []()
* []()
* []()

<p align="right">(<a href="#readme-top">back to top</a>)</p>

[github-cloud-commmis]: https://github.com/alexandre-girault/cloud-commis
[stars-shield]: https://img.shields.io/github/stars/github_username/repo_name.svg?style=for-the-badge
[stars-url]: https://github.com/alexandre-girault/cloud-commis/stargazers
[issues-shield]: https://img.shields.io/github/issues/github_username/repo_name.svg?style=for-the-badge
[issues-url]: https://github.com/alexandre-girault/cloud-commis/issues
[license-shield]: https://img.shields.io/github/license/github_username/repo_name.svg?style=for-the-badge
[license-url]: https://github.com/alexandre-girault/cloud-commis/blob/master/LICENSE.txt
[Golang-logo]: https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Blue.png
[Golang-url]: https://go.dev/
[Bootstrap.com]: https://img.shields.io/badge/Bootstrap-563D7C?style=for-the-badge&logo=bootstrap&logoColor=white
[Bootstrap-url]: https://getbootstrap.com
[JQuery.com]: https://img.shields.io/badge/jQuery-0769AD?style=for-the-badge&logo=jquery&logoColor=white
[JQuery-url]: https://jquery.com 