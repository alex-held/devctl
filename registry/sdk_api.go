package registry

type SDKApi interface {
    GetSDKs() (sdks []string, err error)
}

func (client GithubRegistryApiClient) GetSDKs() (sdks []string, err error)  {
    uri := client.getContentUrl("sdk")
    files, err := client.getFiles(uri)
    if err != nil {
        return sdks, err
    }
    for _, file := range files {
        if file.Type == "dir"{
            sdks = append(sdks, file.Name)
        }
    }

    return sdks, nil
}
