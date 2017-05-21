package com.xchat.stevenzack.langenius;

import android.content.SharedPreferences;
import android.graphics.Color;
import android.os.Environment;
import android.os.Handler;
import android.os.Message;
import android.support.v4.content.ContextCompat;
import android.support.v4.os.EnvironmentCompat;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.util.SparseBooleanArray;
import android.view.Menu;
import android.view.MenuInflater;
import android.view.MenuItem;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.CheckedTextView;
import android.widget.ListView;
import android.widget.RadioButton;
import android.widget.SimpleAdapter;
import android.widget.Toast;

import java.io.File;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Objects;

public class Main2Activity extends AppCompatActivity {
    private Button button;
    private ListView listView;
    private List<String> filenames=new ArrayList<>();
    private String currentPath;
    private boolean selectAllStatus=false;
    private Handler handler=new Handler(){
        @Override
        public void handleMessage(Message msg) {
            switch (msg.arg1){
                case 0://reload listview   ( on currentPath changed )
                    getSupportActionBar().setSubtitle(currentPath);
                    selectAllStatus=false;

                    File[] listOfFiles=new File(currentPath).listFiles();
//                    for (int i=0;i<listView.getCheckedItemCount();i++){
//                        listView.setItemChecked(i,false);
//                    }
                    filenames.clear();
                    if (listOfFiles!=null){
                        for (int i=0;i<listOfFiles.length;i++){
                            if (listOfFiles[i].getName().startsWith("."))
                                continue;
                            if (listOfFiles[i].isFile()){
                                filenames.add(listOfFiles[i].getName());
                            }else if (listOfFiles[i].isDirectory()){
                                filenames.add(0,listOfFiles[i].getName()+"/");
                            }
                        }
                    }
                    ArrayAdapter<String> arrayAdapter = new ArrayAdapter<>(Main2Activity.this, android.R.layout.simple_list_item_multiple_choice, filenames);
                    listView.setAdapter(arrayAdapter);
                    listView.setSelection((int)msg.obj);
                    break;
                case 1://select all items
                    selectAllStatus=!selectAllStatus;
                    for (int i=0;i<listView.getAdapter().getCount();i++){
                        if (new File(currentPath+filenames.get(i)).isFile()){
                            listView.setItemChecked(i,selectAllStatus);
                        }
                    }
                    break;
            }
        }
    },mainHandler;
    private String TAG="===MainActivity=====";

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main2);
        button=(Button)findViewById(R.id.filechooser_button);
        button.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                SparseBooleanArray bs = listView.getCheckedItemPositions();
                for (int i=0;i<listView.getAdapter().getCount();i++){
                    if (bs.get(i)){
                        if (HandlerConverter.handler!=null){
                            Message msg=new Message();msg.arg1=2;msg.obj=currentPath+filenames.get(i);
                            HandlerConverter.handler.sendMessage(msg);
                        }
                    }
                }

                SharedPreferences sp=getSharedPreferences(Main2Activity.this.getString(R.string.sp_settings),MODE_PRIVATE);
                sp.edit().putString("fileChooserDir",currentPath).commit();
                finish();
            }
        });
        listView=(ListView)findViewById(R.id.filechooser_listview);

        SharedPreferences sp=getSharedPreferences(this.getString(R.string.sp_settings),MODE_PRIVATE);
        currentPath=sp.getString("fileChooserDir",Environment.getExternalStorageDirectory().toString()+"/");
        getSupportActionBar().setSubtitle(currentPath);

        listView.setChoiceMode(ListView.CHOICE_MODE_MULTIPLE);
        listView.setOnItemClickListener(new AdapterView.OnItemClickListener() {
            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
                File file=new File(currentPath+filenames.get(position));
                if (file.isDirectory()){
                    currentPath=currentPath+filenames.get(position);
                    SharedPreferences sp=getSharedPreferences("cache",MODE_PRIVATE);
                    sp.edit().putInt("lastVisit",position).commit();
                    Message msg=new Message();msg.arg1=0;msg.obj=0;handler.sendMessage(msg);//reload listview
                }
            }
        });
        ArrayAdapter<String> arrayAdapter = new ArrayAdapter<>(this, android.R.layout.simple_list_item_multiple_choice, filenames);
        listView.setAdapter(arrayAdapter);
        Message msg=new Message();msg.arg1=0;msg.obj=0;handler.sendMessage(msg);//reload listview
    }
    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        MenuInflater inflater = getMenuInflater();
        inflater.inflate(R.menu.filechooser_menu, menu);
        return true;
    }

    @Override
    public void onBackPressed() {
        if (currentPath.length()<=(Environment.getExternalStorageDirectory().toString()+"/").length()){
            SharedPreferences sp=getSharedPreferences(this.getString(R.string.sp_settings),MODE_PRIVATE);
            sp.edit().putString("fileChooserDir",Environment.getExternalStorageDirectory().toString()+"/").commit();
            finish();
        }else{
            String[] strings=currentPath.split("/");
            String newPath="/";
            for (int i=1;i<strings.length-1;i++){
                newPath=newPath+strings[i]+"/";
            }
            currentPath=newPath;
            SharedPreferences sp=getSharedPreferences("cache",MODE_PRIVATE);
            Message msg=new Message();msg.arg1=0;msg.obj=sp.getInt("lastVisit",0);handler.sendMessage(msg);//reload listview
        }
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        switch (item.getItemId()){
            case R.id.filechooser_menu_all:
                Message msg=new Message();msg.arg1=1;handler.sendMessage(msg);
                return true;
            default:
                return true;
        }
    }
}
