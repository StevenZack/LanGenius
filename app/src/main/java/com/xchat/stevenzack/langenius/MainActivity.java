package com.xchat.stevenzack.langenius;

import android.Manifest;
import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.pm.PackageManager;
import android.database.Cursor;
import android.net.Uri;
import android.os.Build;
import android.os.Handler;
import android.os.Message;
import android.provider.MediaStore;
import android.support.design.widget.FloatingActionButton;
import android.support.v4.app.ActivityCompat;
import android.support.v4.content.ContextCompat;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ImageButton;
import android.widget.ListView;
import android.widget.SimpleAdapter;
import android.widget.TextView;
import android.widget.Toast;

import java.io.File;
import java.net.Inet6Address;
import java.net.InetAddress;
import java.net.NetworkInterface;
import java.net.SocketException;
import java.net.URISyntaxException;
import java.util.ArrayList;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.List;

import LanGenius.JavaHandler;
import LanGenius.LanGenius;


public class MainActivity extends AppCompatActivity {
    private ClipboardManager clipboardManager;
    private String TAG="Main";
    private Handler handler=new Handler(){
        @Override
        public void handleMessage(Message msg) {
            switch (msg.arg1){
                case 0://onClipboard received
                    clipboardManager.setPrimaryClip(ClipData.newPlainText("Copied text",msg.obj.toString()));
                    Toast.makeText(MainActivity.this,MainActivity.this.getString(R.string.newClipboard),Toast.LENGTH_SHORT).show();
                    break;
                case 1://on File received
                    Toast.makeText(MainActivity.this,MainActivity.this.getString(R.string.newFile)+msg.obj.toString(),Toast.LENGTH_SHORT).show();
                    break;
                case 2:
                    if (msg.obj!=null) {
                        String path = msg.obj.toString();
                        try {
                            Log.d(TAG, "onActivityResult: PATH ==" + path);
                            String[] strs = path.split("/");
                            HashMap<String,String> hashMap=new HashMap<>();
                            hashMap.put("FileName",strs[strs.length-1]);
                            strings.add(hashMap);
                            simpleAdapter.notifyDataSetChanged();
                            LanGenius.addFile(path);
                        } catch (Exception e) {
                            Log.d(TAG, "onActivityResult: Exception  == " + e.toString());
                            Toast.makeText(MainActivity.this, MainActivity.this.getString(R.string.addFileFailed), Toast.LENGTH_SHORT).show();
                        }
                    }else {
                        Toast.makeText(MainActivity.this, MainActivity.this.getString(R.string.addFileFailed), Toast.LENGTH_SHORT).show();
                    }
                    break;
            }
        }
    };
    private TextView txt_ip;
    private ListView listView;
    private List<HashMap<String,String>> strings=new ArrayList<>();
    private SimpleAdapter simpleAdapter;
    private FloatingActionButton floatingActionButton;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        clipboardManager=(ClipboardManager)getSystemService(CLIPBOARD_SERVICE);
        clipboardManager.addPrimaryClipChangedListener(new ClipboardManager.OnPrimaryClipChangedListener() {
            @Override
            public void onPrimaryClipChanged() {
                LanGenius.setClipboard(clipboardManager.getPrimaryClip().getItemAt(0).getText().toString());
            }
        });
        if (clipboardManager.getPrimaryClip()!=null) {
            LanGenius.setClipboard(clipboardManager.getPrimaryClip().getItemAt(0).getText().toString());
        }
        txt_ip=(TextView)findViewById(R.id.txt_hostname);
        String str_IP=getHostIP();
        txt_ip.setText(MainActivity.this.getString(R.string.websiteAddress)+(str_IP==null?"localhost":str_IP)+":4444");
        isStoragePermissionGranted();
        String lang=MainActivity.this.getString(R.string.language);
        LanGenius.start(lang,new MyJavaHandler());
        ((ImageButton)findViewById(R.id.bt_openbrowser)).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                String tempstr=getHostIP();
                if (tempstr==null)
                    tempstr="localhost";
                Intent browserIntent = new Intent(Intent.ACTION_VIEW, Uri.parse("http://"+tempstr+":4444"));
                startActivity(browserIntent);
            }
        });
        listView=(ListView)findViewById(R.id.main_listview);
//        simpleAdapter=new SimpleAdapter(this,strings,R.layout.listview_item,new String[]{"FileName"},new int[]{R.id.});
        simpleAdapter=new SimpleAdapter(this,strings,R.layout.listview_item,new String[]{"FileName"},new int[]{R.id.listview_txt_filename});
        listView.setAdapter(simpleAdapter);
        floatingActionButton=(FloatingActionButton)findViewById(R.id.floatingActionButton);
        floatingActionButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                showFileChooser();
            }
        });
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP) {
            getWindow().setNavigationBarColor(ContextCompat.getColor(this,R.color.colorPrimaryDark));
        }
        SharedPreferences sp_settings=getSharedPreferences(MainActivity.this.getString(R.string.sp_settings),MODE_PRIVATE);
        String str=sp_settings.getString(this.getString(R.string.sp_sub_frcv_path),this.getString(R.string.storagepath));
        ((TextView)findViewById(R.id.main_frcv_path)).setText(str);
    }
    @Override
    protected void onDestroy() {
        LanGenius.stop();
        super.onDestroy();
    }

    private void showFileChooser() {
        Intent intent = new Intent(Intent.ACTION_GET_CONTENT);
        intent.setType("*/*");
        intent.addCategory(Intent.CATEGORY_OPENABLE);
        try {
            startActivityForResult( Intent.createChooser(intent, "Select a File to Upload"), 2233);
        } catch (android.content.ActivityNotFoundException ex) {
            Toast.makeText(this, "Please install a File Manager.",  Toast.LENGTH_SHORT).show();
        }
    }

    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        switch (requestCode){
            case 2233:
                if (resultCode==RESULT_OK){
                    Uri uri=data.getData();
                    Log.d("spy","##FileSharer: uri="+uri.toString());
//                    try {
//                        String path=getPath(this,uri);
//                        Log.d(TAG, "onActivityResult: PATH ================"+path);
//                        String[] strs = path.split("/");
//                        strings.add(strs[strs.length-1]);
//                        listView.notifyAll();
//                        LanGenius.addFile(path);
//                    } catch (Exception e) {
//                        Log.d(TAG, "onActivityResult: "+e.toString());
//                        Toast.makeText(MainActivity.this,MainActivity.this.getString(R.string.addFileFailed),Toast.LENGTH_SHORT).show();
//                    }
//                    String[] proj = {MediaStore.Images.Media.DATA};
//                    Cursor actualimagecursor = managedQuery(uri, proj, null, null, null);
//                    int actual_image_column_index = actualimagecursor.getColumnIndexOrThrow(MediaStore.Images.Media.DATA);
//                    actualimagecursor.moveToFirst();
//                    String img_path = actualimagecursor.getString(actual_image_column_index);
//                    File file = new File(img_path);
//                    Toast.makeText(MainActivity.this, file.toString(), Toast.LENGTH_SHORT).show();
                    String path=FileUtils.getPath(MainActivity.this,uri);
                    Message msg=new Message();
                    msg.arg1=2;
                    msg.obj=path;
                    handler.sendMessage(msg);
                }
                break;
        }
        super.onActivityResult(requestCode, resultCode, data);
    }

    public  boolean isStoragePermissionGranted() {
        if (Build.VERSION.SDK_INT >= 23) {
            if (checkSelfPermission(Manifest.permission.WRITE_EXTERNAL_STORAGE)
                    == PackageManager.PERMISSION_GRANTED&&
                    checkSelfPermission(Manifest.permission.READ_EXTERNAL_STORAGE)
                            == PackageManager.PERMISSION_GRANTED) {
                Log.v(TAG,"Permission is granted");
                return true;
            } else {

                Log.v(TAG,"Permission is revoked");
                ActivityCompat.requestPermissions(this, new String[]{Manifest.permission.WRITE_EXTERNAL_STORAGE}, 1);
                ActivityCompat.requestPermissions(this, new String[]{Manifest.permission.READ_EXTERNAL_STORAGE}, 1);
                return false;
            }
        }
        else { //permission is automatically granted on sdk<23 upon installation
            Log.v(TAG,"Permission is granted");
            return true;
        }
    }
    class MyJavaHandler implements JavaHandler{
        @Override
        public void onClipboardReceived(String s) {
            Message msg=new Message();
            msg.arg1=0;
            msg.obj=s;
            handler.sendMessage(msg);
        }

        @Override
        public void onFileReceived(String s) {
            if (isStoragePermissionGranted()){
                Message msg=new Message();
                msg.arg1=1;
                msg.obj=s;
                handler.sendMessage(msg);
            }
        }
    }
    public static String getHostIP() {

        String hostIp = null;
        try {
            Enumeration nis = NetworkInterface.getNetworkInterfaces();
            InetAddress ia = null;
            while (nis.hasMoreElements()) {
                NetworkInterface ni = (NetworkInterface) nis.nextElement();
                Enumeration<InetAddress> ias = ni.getInetAddresses();
                while (ias.hasMoreElements()) {
                    ia = ias.nextElement();
                    if (ia instanceof Inet6Address) {
                        continue;// skip ipv6
                    }
                    String ip = ia.getHostAddress();
                    if (!"127.0.0.1".equals(ip)) {
                        hostIp = ia.getHostAddress();
                        break;
                    }
                }
            }
        } catch (SocketException e) {
            Log.i("yao", "SocketException");
            e.printStackTrace();
        }
        return hostIp;

    }

}
