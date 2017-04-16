package com.xchat.stevenzack.langenius;

import android.content.Context;
import android.content.SharedPreferences;
import android.os.Build;
import android.os.Handler;
import android.os.Message;
import android.support.v4.app.Fragment;
import android.support.v4.app.FragmentManager;
import android.support.v4.app.FragmentPagerAdapter;
import android.support.v4.content.ContextCompat;
import android.support.v4.view.ViewPager;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.text.Editable;
import android.text.TextUtils;
import android.text.TextWatcher;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.EditText;
import android.widget.ImageButton;
import android.widget.ImageView;
import android.widget.Spinner;
import android.widget.TextView;
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
                case 0://toast
                    Toast.makeText(KCActivity.this,msg.obj.toString(),Toast.LENGTH_SHORT).show();
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
    private static ArrayAdapter<String> arrayAdapter=null;
    private static long[] currentDevice=new long[]{-1,-1,-1,-1};
    private static Context context;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_kc);
        context=this;
        arrayAdapter=new ArrayAdapter<>(KCActivity.this,android.R.layout.simple_spinner_dropdown_item,list_device);
        ((ViewPager)findViewById(R.id.kc_viewpager)).setAdapter(new SectionsPagerAdapter(getSupportFragmentManager()));
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP) {
            getWindow().setNavigationBarColor(ContextCompat.getColor(this,R.color.colorPrimaryDark));
        }
        LanGenius.startKC(new MyKCHandler());
        String str=LanGenius.getUDPConnections();
        if (!str.equals("")) {
            Log.d("Fuck", "onCreate: =======================str = "+str);
            String[] sd = str.split("#");
            list_device.clear();
            for (int j = 0; j < sd.length; j++) {
                list_device.add(sd[j]);
            }
            arrayAdapter.notifyDataSetChanged();
        }
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
    public static class PlaceholderFragment extends Fragment{
        private static final String ARG_SECTION_NUMBER = "section_number";

        public PlaceholderFragment() {
        }
        public static PlaceholderFragment newInstance(int sectionNumber) {
            PlaceholderFragment fragment = new PlaceholderFragment();
            Bundle args = new Bundle();
            args.putInt(ARG_SECTION_NUMBER, sectionNumber);
            fragment.setArguments(args);
            return fragment;
        }
        @Override
        public View onCreateView(LayoutInflater inflater, ViewGroup container,
                                 Bundle savedInstanceState) {
            View rootView = inflater.inflate(R.layout.kc_frag_main, container, false);
            ImageButton frag_bt=(ImageButton)rootView.findViewById(R.id.kc_frag_press_bt) ;
            Spinner frag_spinner=(Spinner)rootView.findViewById(R.id.kc_frag_spinner);
            final EditText editText=(EditText)rootView.findViewById(R.id.kc_frag_edit_text);
            final SharedPreferences sp_kc=KCActivity.context.getSharedPreferences(KCActivity.context.getString(R.string.sp_settings),MODE_PRIVATE);
            final int i;
            switch (getArguments().getInt(ARG_SECTION_NUMBER)){
                case 1:
                    frag_bt.setImageResource(R.drawable.circle_primary_dark);
                    i=0;
                    break;
                case 2:
                    frag_bt.setImageResource(R.drawable.circle_primary_dark2);
                    i=1;
                    break;
                case 3:
                    frag_bt.setImageResource(R.drawable.circle_primary_dark3);
                    i=2;
                    break;
                case 4:
                    frag_bt.setImageResource(R.drawable.circle_primary_dark4);
                    i=3;
                    break;
            }
            frag_spinner.setAdapter(arrayAdapter);
            frag_spinner.setOnItemSelectedListener(new AdapterView.OnItemSelectedListener() {
                @Override
                public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
                    currentDevice[getArguments().getInt(ARG_SECTION_NUMBER)-1]=position;
                }
                @Override
                public void onNothingSelected(AdapterView<?> parent) {
                }
            });
            editText.setText(sp_kc.getString("keymap"+String.valueOf(getArguments().getInt(ARG_SECTION_NUMBER)-1),""));
            editText.addTextChangedListener(new TextWatcher() {
                @Override
                public void beforeTextChanged(CharSequence s, int start, int count, int after) {
                }
                @Override
                public void onTextChanged(CharSequence s, int start, int before, int count) {
                    sp_kc.edit().putString("keymap"+String.valueOf(getArguments().getInt(ARG_SECTION_NUMBER)-1),s.toString()).commit();
                }
                @Override
                public void afterTextChanged(Editable s) {
                }
            });
            frag_bt.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    String string=editText.getText().toString();
                    string.trim();
                    Log.d("Fuck", "onClick: You just clicked bt "+String.valueOf(getArguments().getInt(ARG_SECTION_NUMBER)-1)+" string ="+string);
                    if (!string.equals("")&&currentDevice[getArguments().getInt(ARG_SECTION_NUMBER)-1]>-1){
                        String[] sts=string.split(" ");
                        String[] sts_r=new String[sts.length];
                        for (int n=0;n<sts.length;n++){
                            sts_r[sts.length-1-n]=sts[n];
                        }
                        Log.d("Fragment=====", "onClick: "+TextUtils.join("#",sts_r));
                        try {
                            LanGenius.sendKC(TextUtils.join("#",sts_r),currentDevice[getArguments().getInt(ARG_SECTION_NUMBER)-1]);
                        }catch (Exception e){
                            Log.d("Fragment=========", "onClick: "+e.toString());
                        }
                    }
                }
            });
            return rootView;
        }
    }
    public class SectionsPagerAdapter extends FragmentPagerAdapter {

        public SectionsPagerAdapter(FragmentManager fm) {
            super(fm);
        }

        @Override
        public Fragment getItem(int position) {
            return PlaceholderFragment.newInstance(position + 1);
        }

        @Override
        public int getCount() {
            // Show 3 total pages.
            return 4;
        }

        @Override
        public CharSequence getPageTitle(int position) {
            switch (position) {
                case 0:
                    return "SECTION 1";
                case 1:
                    return "SECTION 2";
                case 2:
                    return "SECTION 3";
                case 3:
                    return "SECTION 4";
            }
            return null;
        }
    }
}
