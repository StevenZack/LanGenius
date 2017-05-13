package com.xchat.stevenzack.langenius;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ListView;
import android.widget.RadioButton;
import android.widget.SimpleAdapter;
import android.widget.Toast;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Objects;

public class Main2Activity extends AppCompatActivity {
    private Button button;
    private ListView listView;
    private List<String> strs=new ArrayList<>();
    private ArrayAdapter arrayAdapter;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main2);
        button=(Button)findViewById(R.id.filechooser_button);
        button.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Toast.makeText(Main2Activity.this,"clicked button",Toast.LENGTH_SHORT).show();
            }
        });

        listView=(ListView)findViewById(R.id.filechooser_listview);

        strs.add("asdda.jpg");
        strs.add("asdda.jpg");
        strs.add("asdda.jpg");
        arrayAdapter=new ArrayAdapter<>(this,android.R.layout.simple_list_item_multiple_choice,strs);
        listView.setChoiceMode(ListView.CHOICE_MODE_MULTIPLE);
        listView.setOnItemClickListener(new AdapterView.OnItemClickListener() {
            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
                Toast.makeText(Main2Activity.this,String.valueOf(listView.getCheckedItemCount()),Toast.LENGTH_SHORT).show();
            }
        });
        listView.setAdapter(arrayAdapter);
    }
}
