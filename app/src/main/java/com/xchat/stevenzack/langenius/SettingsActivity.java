package com.xchat.stevenzack.langenius;

import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.content.SharedPreferences;
import android.net.Uri;
import android.os.Build;
import android.os.Environment;
import android.os.Handler;
import android.os.Message;
import android.support.v4.content.ContextCompat;
import android.support.v7.app.AlertDialog;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.text.InputType;
import android.util.Log;
import android.view.View;
import android.view.inputmethod.InputMethodManager;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Toast;

import java.io.File;

import LanGenius.LanGenius;

public class SettingsActivity extends AppCompatActivity {
    private Button bt_filercvpath,bt_default_port,bt_html;
    private Handler handler=new Handler(){
        @Override
        public void handleMessage(Message msg) {
            switch (msg.arg1){
                case 0:
                    break;
            }
        }
    };
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_settings);
        final SharedPreferences sp_settings=getSharedPreferences(SettingsActivity.this.getString(R.string.sp_settings),MODE_PRIVATE);
        bt_filercvpath=(Button)findViewById(R.id.set_bt_frcvpath);
        bt_filercvpath.setText(this.getString(R.string.storagepath)+sp_settings.getString(this.getString(R.string.sp_sub_frcv_path), Environment.getExternalStorageDirectory().toString()+"/"));
        bt_filercvpath.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                AlertDialog.Builder builder = new AlertDialog.Builder(SettingsActivity.this);
                builder.setTitle(SettingsActivity.this.getString(R.string.str_setDefaultFileRecvPath));
                final EditText input = new EditText(SettingsActivity.this);
                input.setInputType(InputType.TYPE_CLASS_TEXT);
                input.setText(sp_settings.getString(SettingsActivity.this.getString(R.string.sp_sub_frcv_path),Environment.getExternalStorageDirectory().toString()+"/"));
                builder.setView(input);
                builder.setPositiveButton(SettingsActivity.this.getString(R.string.str_ok), new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        String m_Text = input.getText().toString();
                        if (!m_Text.endsWith("/")){
                            m_Text=m_Text+"/";
                        }
                        sp_settings.edit().putString(SettingsActivity.this.getString(R.string.sp_sub_frcv_path),m_Text).commit();
                        bt_filercvpath.setText(SettingsActivity.this.getString(R.string.storagepath)+sp_settings.getString(SettingsActivity.this.getString(R.string.sp_sub_frcv_path),Environment.getExternalStorageDirectory().toString()+"/"));
                        LanGenius.setStoragePath(m_Text);
                    }
                });
                builder.setNegativeButton(SettingsActivity.this.getString(R.string.str_cancel), new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        dialog.cancel();
                    }
                });

                builder.show();
            }
        });
        bt_default_port=(Button)findViewById(R.id.set_bt_default_port);
        bt_default_port.setText(SettingsActivity.this.getString(R.string.str_default_port)+sp_settings.getString(this.getString(R.string.sp_sub_port),this.getString(R.string.default_port)));
        bt_default_port.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                AlertDialog.Builder builder = new AlertDialog.Builder(SettingsActivity.this);
                builder.setTitle(SettingsActivity.this.getString(R.string.str_setDefaultPort));
                final EditText input = new EditText(SettingsActivity.this);
                input.setInputType(InputType.TYPE_CLASS_NUMBER);
                input.setHint("1000~65535");
                builder.setView(input);
                builder.setPositiveButton(SettingsActivity.this.getString(R.string.str_ok), new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        if (input.getText()==null){
                            Toast.makeText(SettingsActivity.this, SettingsActivity.this.getString(R.string.str_must_between), Toast.LENGTH_SHORT).show();
                        }else {
                            final String tmpstr = input.getText().toString();
                            if (tmpstr == null || tmpstr.equals("") || Integer.valueOf(tmpstr) < 1000 || Integer.valueOf(tmpstr) > 65535) {
                                Toast.makeText(SettingsActivity.this, SettingsActivity.this.getString(R.string.str_must_between), Toast.LENGTH_SHORT).show();
                            } else {
                                String m_Text = ":" + input.getText().toString();
                                sp_settings.edit().putString(SettingsActivity.this.getString(R.string.sp_sub_port), m_Text).commit();
                                bt_default_port.setText(SettingsActivity.this.getString(R.string.str_default_port) + sp_settings.getString(SettingsActivity.this.getString(R.string.sp_sub_port), SettingsActivity.this.getString(R.string.default_port)));
                            }
                        }
                    }
                });
                builder.setNegativeButton(SettingsActivity.this.getString(R.string.str_cancel), new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        dialog.cancel();
                    }
                });

                builder.show();
            }
        });
        bt_html=(Button)findViewById(R.id.set_bt_html);
        bt_html.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                AlertDialog.Builder builder = new AlertDialog.Builder(SettingsActivity.this);
                builder.setTitle(SettingsActivity.this.getString(R.string.str_set_port));
                final EditText input = new EditText(SettingsActivity.this);
                input.setInputType(InputType.TYPE_CLASS_NUMBER);
                input.setHint("1000~65535");
                input.setText(sp_settings.getString("StaticSitePort",""));

                builder.setView(input);
                builder.setPositiveButton(SettingsActivity.this.getString(R.string.str_ok), new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        if (input.getText()==null){
                            Toast.makeText(SettingsActivity.this, SettingsActivity.this.getString(R.string.str_must_between), Toast.LENGTH_SHORT).show();
                        }else {
                            final String tmpstr = input.getText().toString();
                            if (tmpstr == null || tmpstr.equals("")|| Integer.valueOf(tmpstr) < 1000 || Integer.valueOf(tmpstr) > 65535 ) {
                                Toast.makeText(SettingsActivity.this, SettingsActivity.this.getString(R.string.str_must_between), Toast.LENGTH_SHORT).show();
                            } else if ((":" + tmpstr).equals(sp_settings.getString(SettingsActivity.this.getString(R.string.sp_sub_port), ""))) {
                                Toast.makeText(SettingsActivity.this, SettingsActivity.this.getString(R.string.str_cannot_use_LanGenius_port), Toast.LENGTH_LONG).show();
                            } else {
                                AlertDialog.Builder builder = new AlertDialog.Builder(SettingsActivity.this);
                                builder.setTitle(SettingsActivity.this.getString(R.string.str_select_dir));
                                final EditText input = new EditText(SettingsActivity.this);
                                input.setInputType(InputType.TYPE_CLASS_TEXT);
                                input.setText(sp_settings.getString("staticSiteDir", Environment.getExternalStorageDirectory() + "/"));
                                builder.setView(input);
                                builder.setPositiveButton(SettingsActivity.this.getString(R.string.str_ok), new DialogInterface.OnClickListener() {
                                    @Override
                                    public void onClick(DialogInterface dialog, int which) {
                                        if (input.getText()==null){
                                            Toast.makeText(SettingsActivity.this, SettingsActivity.this.getString(R.string.str_path_incorrect), Toast.LENGTH_SHORT).show();
                                        }else {
                                            String tmpPath = input.getText().toString();
                                            if (tmpPath == null || !tmpPath.startsWith("/") || !(new File(tmpPath).isDirectory())) {
                                                Toast.makeText(SettingsActivity.this, SettingsActivity.this.getString(R.string.str_path_incorrect), Toast.LENGTH_SHORT).show();
                                            } else {
                                                LanGenius.startStaticSite(":" + tmpstr, tmpPath);
                                                Log.d("TAG", "onClick: tmpPath"+tmpPath);
                                                sp_settings.edit().putString("StaticSitePort", tmpstr).putString("staticSiteDir", tmpPath).commit();
                                                bt_html.setText(SettingsActivity.this.getString(R.string.str_running_on) + LanGenius.getIP() + ":" + tmpstr);
                                                bt_html.setEnabled(false);
                                            }
                                        }
                                    }
                                });
                                builder.setNegativeButton(SettingsActivity.this.getString(R.string.str_cancel), new DialogInterface.OnClickListener() {
                                    @Override
                                    public void onClick(DialogInterface dialog, int which) {
                                        dialog.cancel();
                                    }
                                });

                                builder.show();
                            }
                        }
                    }
                });
                builder.setNegativeButton(SettingsActivity.this.getString(R.string.str_cancel), new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        dialog.cancel();
                    }
                });

                builder.show();
            }
        });
        if (LanGenius.isStaticSiteRunning()){
            bt_html.setText(SettingsActivity.this.getString(R.string.str_running_on)+LanGenius.getIP()+":"+sp_settings.getString("StaticSitePort","NULL"));
            bt_html.setEnabled(false);
        }
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP) {
            getWindow().setNavigationBarColor(ContextCompat.getColor(this,R.color.colorPrimaryDark));
        }
    }

}
