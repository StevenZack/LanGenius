package com.xchat.stevenzack.langenius;

import android.os.Handler;
import android.os.Message;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.text.TextUtils;
import android.util.Log;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.EditText;
import android.widget.ImageButton;
import android.widget.Spinner;
import android.widget.Toast;

import java.util.ArrayList;
import java.util.List;

import LanGenius.JavaKCHandler;
import LanGenius.LanGenius;

public class KCActivity extends AppCompatActivity {
    private Handler handler=new Handler(){
        @Override
        public void handleMessage(Message msg) {
            switch (msg.arg1){
                case 0:
                    break;
                case 1://remote device detected
                    list_device.add(msg.obj.toString());
                    arrayAdapter.notifyDataSetChanged();
                    Toast.makeText(KCActivity.this,KCActivity.this.getString(R.string.str_new_device_detected),Toast.LENGTH_SHORT).show();
                    break;
            }
        }
    };
    private List<String> list_device=new ArrayList<>();
    private ArrayAdapter<String> arrayAdapter=null;
    private Spinner spinner;
    private long currentDevice=-1;
    private ImageButton imageButton;
    private EditText edittext;
    private String TAG="===KC-activity=======";

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_kc);
        arrayAdapter=new ArrayAdapter<>(KCActivity.this,android.R.layout.simple_spinner_dropdown_item,list_device);
        spinner=(Spinner)findViewById(R.id.kc_spinner);
        spinner.setAdapter(arrayAdapter);
        spinner.setOnItemSelectedListener(new AdapterView.OnItemSelectedListener() {
            @Override
            public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
                currentDevice=position;
            }
            @Override
            public void onNothingSelected(AdapterView<?> parent) {
            }
        });
        edittext=(EditText)findViewById(R.id.kc_edit_text);
        imageButton=(ImageButton)findViewById(R.id.kc_press_bt);
        imageButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                String string=edittext.getText().toString();
                string.trim();
                if (!string.equals("")&&currentDevice>-1){
                    String[] sts=string.split(" ");
                    String[] sts_reversed=new String[sts.length];
                    for (int i=0;i<sts.length;i++){
                        sts_reversed[sts.length-1-i]=sts[i];
                    }
                    Log.d(TAG, "onClick: "+ TextUtils.join("#",sts_reversed));
                    try {
                        LanGenius.sendKC(TextUtils.join("#",sts_reversed),currentDevice);
                    } catch (Exception e) {
                        Log.d(TAG, "onClick: "+e.toString());
                        Toast.makeText(KCActivity.this,e.toString(),Toast.LENGTH_SHORT).show();
                    }
                }
            }
        });
        LanGenius.startKC(new MyKCHandler());

    }

    @Override
    protected void onDestroy() {
        LanGenius.stopKC();
        super.onDestroy();
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
