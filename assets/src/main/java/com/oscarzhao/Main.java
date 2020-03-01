package com.oscarzhao;

import java.util.Map;
import java.util.HashMap;

public class Main {
  public static void main(String[] args) {
    Long a = 1L;
    Long b = a;
    System.out.printf("a = %d, b = %d\n", a, b);
    System.out.printf("equal = %b\n", a.equals(b));
    System.out.printf("== = %b\n", a==b);

    System.out.println("------ after -------");
    b = 2L;
    System.out.printf("a = %d, b = %d\n", a, b);
    System.out.printf("equal = %b\n", a.equals(b));
    System.out.printf("== = %b\n", a==b);

    System.out.println("\n\n");
    Map<String, String> ma = new HashMap<String, String>();
    ma.put("a", "b");
    Map<String, String> mb = ma;
    System.out.printf("ma['a']=%s, ma['b']=%s\n", ma.get("a"), ma.get("b"));
    System.out.printf("mb['a']=%s, mb['b']=%s\n", mb.get("a"), mb.get("b"));


    System.out.println("------ after ------ \n");

    mb.put("b", "c");
    System.out.printf("ma['a']=%s, ma['b']=%s\n", ma.get("a"), ma.get("b"));
    System.out.printf("mb['a']=%s, mb['b']=%s\n", mb.get("a"), mb.get("b"));

    System.out.println("\n\n");
    Map<String, String> mc = new HashMap<String, String>();
    mc.put("a", "b");
    mc.put("b", "c");


    System.out.printf("equal = %b, %b\n", ma.equals(mb), ma.equals(mc));
    System.out.printf("== = %b, %b\n", ma == mb, ma == mc);
  }
}
