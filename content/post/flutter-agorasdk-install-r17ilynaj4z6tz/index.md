+++
date = '2023-09-11T21:20:13+09:00'
draft = false
title = 'Flutter Agora SDKの導入方法'
categories = ['external']
tags = []
externalUrl = ''
+++


AgoraSDKの導入を夏休みにどうしてもできなくて
心残りがあったので、今回導入方法をまとめることにしました🔥

# 導入
## アカウント登録
[Agora.io](https://www.agora.io)でアカウント登録します。

https://www.agora.io

## Agora.ioの設定
Agora.ioにログインし、横のバーにあるところから`Project Management`をクリックします。
![image.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/2017749/361f6407-499e-689e-c17c-a4f0ae0e4718.png)

すると、以下のような画面になります。
![image.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/2017749/902c9cd5-8f67-2467-8d88-9f68880e456c.png)

↑の画面で、`Create`を押します。
すると、以下のようなモーダルが出てくるので
わかりやすいプロジェクト名、用途（Use Case）、`Testing mode`としましょう。
![image.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/2017749/124bd810-7be9-8392-e369-dbe4480dd4c4.png)
そして、Submitを押すと、以下のような画面になるので

後々使うApp IDをどこかにメモ（コピー）しましょう。
![image.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/2017749/04f726da-d403-6265-4700-f5cda02813e7.png)

## sample Projectの作成
プロジェクトを作成します。
IntelliJなどをお使いの方は、そちらから作成していただいても構いません。
```zsh:ターミナル
flutter create agora_sample
```

## pubspec.yamlの編集
```yaml:pubspec.yaml
dependencies:
# ~~~他のライブラリ等の記述は省略~~~
  agora_rtc_engine: ^6.0.0
  permission_handler: ^9.2.0
# ~~~他のライブラリ等の記述は省略~~~
```

`agora_rtc_engine`のバージョンは以下のURLで
最新のものを確認し、記述しましょう！

https://pub.dev/packages/agora_rtc_engine

## プロジェクトにインストール

```zsh:ターミナル
flutter pub get
```

# 実装
サンプルの実装パートは以下のドキュメントを参考にしています。

https://docs.agora.io/en/video-calling/get-started/get-started-sdk?platform=flutter

## main.dartの編集

### import
```dart:lib/main.dart
import 'dart:async';
import 'package:permission_handler/permission_handler.dart';
import 'package:agora_rtc_engine/agora_rtc_engine.dart';
```

### アプリ自体の編集
`main.dart`の内容を以下のものに差し替えましょう。
（既にあるプロダクトで使用する場合は、カスタムして使いましょう。）
```dart:lib/main.dart
const String appId = "先ほどコピーしたApp IDを入れます";

void main() => runApp(const MaterialApp(home: MyApp()));

class MyApp extends StatefulWidget {
const MyApp({Key? key}) : super(key: key);

@override
_MyAppState createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
    String channelName = "<--Insert channel name here-->";
    String token = "<--Insert authentication token here-->";
    
    int uid = 0; // uid of the local user

    int? _remoteUid; // uid of the remote user
    bool _isJoined = false; // Indicates if the local user has joined the channel
    late RtcEngine agoraEngine; // Agora engine instance

    final GlobalKey<ScaffoldMessengerState> scaffoldMessengerKey
    = GlobalKey<ScaffoldMessengerState>(); // Global key to access the scaffold

    showMessage(String message) {
        scaffoldMessengerKey.currentState?.showSnackBar(SnackBar(
        content: Text(message),
        ));
    }
}

```

### 
```dart:lib/main.dart
// ~~~ final GlobalKey<ScaffoldMessengerState> scaffoldMessengerKeyの次の行
    // Build UI
    @override
    Widget build(BuildContext context) {
        return MaterialApp(
        scaffoldMessengerKey: scaffoldMessengerKey,
        home: Scaffold(
            appBar: AppBar(
                title: const Text('Get started with Video Calling'),
            ),
            body: ListView(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
                children: [
                // Container for the local video
                Container(
                    height: 240,
                    decoration: BoxDecoration(border: Border.all()),
                    child: Center(child: _localPreview()),
                ),
                const SizedBox(height: 10),
                //Container for the Remote video
                Container(
                    height: 240,
                    decoration: BoxDecoration(border: Border.all()),
                    child: Center(child: _remoteVideo()),
                ),
                // Button Row
                Row(
                    children: <Widget>[
                    Expanded(
                        child: ElevatedButton(
                            onPressed: _isJoined ? null : () => {join()},
                            child: const Text("Join"),
                        ),
                    ),
                    const SizedBox(width: 10),
                    Expanded(
                        child: ElevatedButton(
                            onPressed: _isJoined ? () => {leave()} : null,
                            child: const Text("Leave"),
                        ),
                    ),
                    ],
                ),
                // Button Row ends
                ],
            )),
        );
    }
    
    // Display local video preview
    Widget _localPreview() {
        if (_isJoined) {
        return AgoraVideoView(
            controller: VideoViewController(
            rtcEngine: agoraEngine,
            canvas: VideoCanvas(uid: uid),
            ),
        );
        } else {
        return const Text(
            'Join a channel',
            textAlign: TextAlign.center,
        );
        }
    }
    
    // Display remote user's video
    Widget _remoteVideo() {
        if (_remoteUid != null) {
        return AgoraVideoView(
            controller: VideoViewController.remote(
            rtcEngine: agoraEngine,
            canvas: VideoCanvas(uid: _remoteUid),
            connection: RtcConnection(channelId: channelName),
            ),
        );
        } else {
            String msg = '';
            if (_isJoined) msg = 'Waiting for a remote user to join';
            return Text(
            msg,
            textAlign: TextAlign.center,
        );
        }
    }

    showMessage(String message) {
        scaffoldMessengerKey.currentState?.showSnackBar(SnackBar(
        content: Text(message),
        ));
    }
}

```


```dart:lib/main.dart
    @override
    void initState() {
        super.initState();
        // Set up an instance of Agora engine
        setupVideoSDKEngine();
    }

    showMessage(String message) {
        scaffoldMessengerKey.currentState?.showSnackBar(SnackBar(
        content: Text(message),
        ));
    }
}
```

```dart:lib/main.dart
    Future<void> setupVideoSDKEngine() async {
        // retrieve or request camera and microphone permissions
        await [Permission.microphone, Permission.camera].request();
    
        //create an instance of the Agora engine
        agoraEngine = createAgoraRtcEngine();
        await agoraEngine.initialize(const RtcEngineContext(
        appId: appId
        ));
    
        await agoraEngine.enableVideo();
        
        // Register the event handler
        agoraEngine.registerEventHandler(
        RtcEngineEventHandler(
            onJoinChannelSuccess: (RtcConnection connection, int elapsed) {
            showMessage("Local user uid:${connection.localUid} joined the channel");
            setState(() {
                _isJoined = true;
            });
            },
            onUserJoined: (RtcConnection connection, int remoteUid, int elapsed) {
            showMessage("Remote user uid:$remoteUid joined the channel");
            setState(() {
                _remoteUid = remoteUid;
            });
            },
            onUserOffline: (RtcConnection connection, int remoteUid,
                UserOfflineReasonType reason) {
            showMessage("Remote user uid:$remoteUid left the channel");
            setState(() {
                _remoteUid = null;
            });
            },
        ),
        );
    }

    showMessage(String message) {
        scaffoldMessengerKey.currentState?.showSnackBar(SnackBar(
        content: Text(message),
        ));
    }
}
```

```dart:lib/main.dart
    void  join() async {
        await agoraEngine.startPreview();
    
        // Set channel options including the client role and channel profile
        ChannelMediaOptions options = const ChannelMediaOptions(
            clientRoleType: ClientRoleType.clientRoleBroadcaster,
            channelProfile: ChannelProfileType.channelProfileCommunication,
        );
    
        await agoraEngine.joinChannel(
            token: token,
            channelId: channelName,
            options: options,
            uid: uid,
        );
    }
    void leave() {
        setState(() {
        _isJoined = false;
        _remoteUid = null;
        });
        agoraEngine.leaveChannel();
    }
    // Clean up the resources when you leave
    @override
    void dispose() async {
        await agoraEngine.leaveChannel();
        super.dispose();
    }

    showMessage(String message) {
        scaffoldMessengerKey.currentState?.showSnackBar(SnackBar(
        content: Text(message),
        ));
    }
}
```
