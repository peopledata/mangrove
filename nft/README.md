
# 如何使用 go 来部署合约

## 安装 solc

可以直接从 github release 页面下载，比如我这里是 Mac m1，下载 `solc-macos` 包，然后放入 PATH 路径：

```bash
# 注意下载的 solc 版本和你的合约代码版本保持一致
$ sudo mv ~/Desktop/solc-macos /usr/local/bin/solc
$ chmod +x /usr/local/bin/solc
```

我们这里的合约代码如下所示：

```solidity
//SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
//import "./contracts/token/ERC721/ERC721.sol";
//import "./contracts/utils/Counters.sol";

contract Minty is ERC721 {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIds;

    constructor(string memory tokenName, string memory symbol) ERC721(tokenName, symbol) {
        _setBaseURI("ipfs://");
    }

    function mintToken(address owner, string memory metadataURI)
    public
    returns (uint256)
    {
        _tokenIds.increment();

        uint256 id = _tokenIds.current();
        _safeMint(owner, id);
        _setTokenURI(id, metadataURI);

        return id;
    }
}
```

需要注意我们这里的合约代码依赖 `@openzeppelin` 库，但是当我们使用 `solc` 进行编译的时候会找不到该库文件，出现类似 `Source "@openzeppelin/contracts/utils/Counters.sol" not found: File outside of allowed directories.` 这样的错误。

这个时候我们可以将 `@openzeppelin` 库下面的 `contracts` 文件夹拷贝到合约代码 `Minty.sol` 同目录下面，然后将依赖引用改为相对路径即可。


## 编译

接下来我们可以通过上面安装的 solidity 编译器编译合约并创建合约的 ABI 文件，通过以下命令在 build 目录中创建的。

```shell
$ solc --optimize --abi Minty.sol -o build
```

然后通过 solidity 编译器在 build 目录中通过下面的命令为同一合约创建 BIN 文件。

```shell
$ solc --optimize --bin ./Minty.sol -o build
```

> 现在我们在 build 目录中同时有合约的 `.abi` 和 `.bin` 文件。

接下来我们需要创建一个 go 文件，我们可以通过 golang 程序与之进行一些对话，因此我们将使用 `.abi` 和 `.bin` 文件，通过 `abigen` 创建一个 go 文件，通过以下命令在我们的 api 目录中创建它。

```shell
$ mkdir api
$ abigen --abi=./build/Minty.abi --bin=./build/Minty.bin --pkg=api --out=./api/Minty.go
```

在 api 文件夹中创建的文件是由 binding 生成的，我们不能编辑它，但在该文件中，我们会发现在 golang 结构体和可交互模式下合约的所有函数和成员变量，代码都是自动生成的，然后我们就可以直接去调用该文件来操作合约了。

> 参考文档：https://medium.com/nerd-for-tech/smart-contract-with-golang-d208c92848a9
