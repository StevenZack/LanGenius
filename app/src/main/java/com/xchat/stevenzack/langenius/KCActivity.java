package com.xchat.stevenzack.langenius;

import android.os.Handler;
import android.os.Message;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.widget.ArrayAdapter;
import android.widget.Spinner;

import java.util.ArrayList;
import java.util.List;

import LanGenius.JavaKCHandler;

public class KCActivity extends AppCompatActivity {
    private Handler handler=new Handler(){
        @Override
        public void handleMessage(Message msg) {
            switch (msg.arg1){
                case 0:
                    break;
                case 1://remote device detected
                    break;
            }
        }
    };
    private List<String> list_device=new ArrayList<>();
    private ArrayAdapter<String> arrayAdapter=null;
    private Spinner spinner;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_kc);
//        list_device.add("192.168.1.100");
//        list_device.add("192.168.1.101");
        arrayAdapter=new ArrayAdapter<>(KCActivity.this,android.R.layout.simple_spinner_dropdown_item,list_device);
        spinner=(Spinner)findViewById(R.id.kc_spinner);
        spinner.setAdapter(arrayAdapter);
    }
    class MyKCHandler implements JavaKCHandler{
        @Override
        public void onDeviceDetected(String s) {
            Message msg=new Message();
            msg.arg1=1;
            msg.obj=s;
            handler.sendMessage(msg);
        }
    }
}
